// Package lockcheck checks for lock misuse.
package lockcheck

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/cfg"
)

// Analyzer is the lockcheck analyzer.
var Analyzer = &analysis.Analyzer{
	Name: "lockcheck",
	Doc:  "reports methods that violate locking conventions",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		ctrlflow.Analyzer,
	},
}

// checkLockSafety is the main logic function for lockcheck
func checkLockSafety(pass *analysis.Pass, fd *ast.FuncDecl, recv, recvMu types.Object) {
	name := fd.Name.String()
	recvIsPrivileged := managesOwnLocking(name)

	cfgs := pass.ResultOf[ctrlflow.Analyzer].(*ctrlflow.CFGs)

	// isLock is a helper that checks for mu.Lock()
	isLock := func(block ast.Node) bool {
		var found bool
		ast.Inspect(block, func(n ast.Node) bool {
			if found {
				return false
			}
			found = isMutexCall(pass, recvMu, n, "Lock")
			// don't descend into DeferStmt or FuncLit
			if _, ok := n.(*ast.DeferStmt); ok {
				return false
			} else if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			return true
		})
		return found
	}
	// isUnlock is a helper that checks for mu.Unlock()
	isUnlock := func(block ast.Node) bool {
		var found bool
		ast.Inspect(block, func(n ast.Node) bool {
			if found {
				return false
			}
			// don't descend into DeferStmt or FuncLit
			if _, ok := n.(*ast.DeferStmt); ok {
				return false
			} else if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			found = isMutexCall(pass, recvMu, n, "Unlock")
			return true
		})
		return found
	}
	// isFieldAccess is a helper that checks for struct field accesses
	isFieldAccess := func(block ast.Node) (field *ast.Ident, ok bool) {
		ast.Inspect(block, func(n ast.Node) bool {
			if field != nil {
				return false // already found
			}
			if _, ok := n.(*ast.FuncLit); ok {
				return false // don't descend into FuncLits
			}
			se, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if x, ok := se.X.(*ast.Ident); ok && pass.TypesInfo.Uses[x] == recv {
				// "sync objects" such as mutexes and threadgroups can be accessed without a lock
				if isSyncObject(pass.TypesInfo.TypeOf(se.Sel)) {
					return true
				}
				field = se.Sel
				return false // no need to search further
			}
			return true
		})
		return field, field != nil
	}
	// isRecvMethodCall is a helper that checks for a method call on
	// a struct/object
	isRecvMethodCall := func(block ast.Node) (method string, ok bool) {
		ast.Inspect(block, func(n ast.Node) bool {
			if method != "" {
				return false // already found
			}
			if _, ok := n.(*ast.FuncLit); ok {
				return false // don't descend into FuncLits
			}
			if ce, ok := n.(*ast.CallExpr); ok {
				if se, ok := ce.Fun.(*ast.SelectorExpr); ok {
					if x, ok := se.X.(*ast.Ident); ok && pass.TypesInfo.Uses[x] == recv {
						method = se.Sel.Name
						return false // no need to search further
					}
				}
			}
			return true
		})
		return method, method != ""
	}
	// isFuncLitCall is a helper that checks for a function literal call
	isFuncLitCall := func(block ast.Node) (litBlock *cfg.Block, ok bool) {
		ast.Inspect(block, func(n ast.Node) bool {
			if litBlock != nil {
				return false // already found
			}
			if ce, ok := n.(*ast.CallExpr); ok {
				switch fn := ce.Fun.(type) {
				case *ast.FuncLit:
					litBlock = cfgs.FuncLit(fn).Blocks[0]
				case *ast.Ident:
					// TODO: make this more generic (currently only handles single assignment)
					if fn.Obj != nil {
						if as, ok := fn.Obj.Decl.(*ast.AssignStmt); ok {
							if lit, ok := as.Rhs[0].(*ast.FuncLit); ok {
								litBlock = cfgs.FuncLit(lit).Blocks[0]
							}
						}
					}
				}
			}
			return true
		})
		return litBlock, litBlock != nil
	}

	// Recursively visit each path through the function, noting the possible
	// lock states at each block.
	type edge struct {
		to, from int32
		locked   bool
	}
	visited := make(map[edge]struct{})
	var checkPath func(*cfg.Block, bool)

	// checkPath is a helper for checking a path for a function
	checkPath = func(b *cfg.Block, lockHeld bool) {
		for _, n := range b.Nodes {
			// Check paths of function literal calls
			if litBlock, ok := isFuncLitCall(n); ok {
				checkPath(litBlock, lockHeld)
				continue
			}
			if isLock(n) {
				// mu.Lock call found
				if !recvIsPrivileged {
					pass.Reportf(n.Pos(), "unprivileged method %s locks mutex", name)
				}
				lockHeld = true
			} else if isUnlock(n) {
				// mu.Unlock call found
				lockHeld = false
			} else if method, ok := isRecvMethodCall(n); ok && !firstWordIs(method, "static") {
				// Method call found that is not a static method
				if recvIsPrivileged {
					// The original object is a managed method
					//
					// First check for calling another managed method while holding
					// a lock.  Ignore threaded methods as those should be called in a go
					// routine and therefore should not create a dead lock
					//
					// TODO: probably should try and add check for calling threaded
					// methods without go routine
					//
					// Second check if we calling an unmanaged method without the lock held
					if managesOwnLocking(method) && !firstWordIs(method, "threaded") && lockHeld {
						pass.Reportf(n.Pos(), "privileged method %s calls privileged method %s while holding mutex", name, method)
					} else if !managesOwnLocking(method) && !lockHeld {
						pass.Reportf(n.Pos(), "privileged method %s calls unprivileged method %s without holding mutex", name, method)
					}
				} else if managesOwnLocking(method) {
					// The original object is not a managed method, so we should not be
					// calling a managed method.
					pass.Reportf(n.Pos(), "unprivileged method %s calls privileged method %s", name, method)
				}
			} else if field, ok := isFieldAccess(n); ok && !isStaticField(field.Name) && !lockHeld {
				// Struct field access found that should be managed by a mutex while no
				// lock is being held
				//
				// NOTE: a method call is also considered a field access, so
				// it's important that we only examine field accesses that
				// aren't method calls (on recv).
				if recvIsPrivileged {
					pass.Reportf(n.Pos(), "privileged method %s accesses %s without holding mutex", name, field)
				}
			}
		}

		for _, succ := range b.Succs {
			e := edge{b.Index, succ.Index, lockHeld}
			if _, ok := visited[e]; ok {
				continue
			}
			visited[e] = struct{}{}
			checkPath(succ, lockHeld)
		}
	}
	checkPath(cfgs.FuncDecl(fd).Blocks[0], false)
}

// containsMutex is a helper that checks if an object contains a mutex
func containsMutex(pass *analysis.Pass, recv types.Object) (types.Object, bool) {
	// Grab the underlying Type Object??
	if p, ok := recv.Type().Underlying().(*types.Pointer); ok {
		// Crab the struct of the pointer's underlying element
		if s, ok := p.Elem().Underlying().(*types.Struct); ok {
			// Iterate over the struct's fields
			for i := 0; i < s.NumFields(); i++ {
				// Check if the field is a Mutex Type and the name is `mu`
				if f := s.Field(i); isMutexType(f.Type()) && s.Field(i).Name() == "mu" {
					return s.Field(i), true
				}
			}
		}
	}
	return nil, false
}

// firstWordIs returns true if name begins with prefix, followed by an uppercase
// letter. For example, firstWordIs("startsUpper", "starts") == true, but
// firstWordIs("starts", "starts") == false.
func firstWordIs(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	suffix := strings.TrimPrefix(name, prefix)
	return len(suffix) > 0 && ast.IsExported(suffix)
}

// isManagedExported returns whether or not a method is a managed exported
// method
func isManagedExported(name string) bool {
	return ast.IsExported(name) && !firstWordIs(name, "Unmanaged")
}

// isMutexCall is a helper that checks for a mutex call
func isMutexCall(pass *analysis.Pass, recvMu types.Object, n ast.Node, muStr string) bool {
	// Check if the Node is an expression followed by an argument list.
	if ce, ok := n.(*ast.CallExpr); ok {
		// Check if the Node is an expression followed by a selector.
		if fnse, ok := ce.Fun.(*ast.SelectorExpr); ok {
			// Check if the selector has the expected suffix.
			if strings.HasSuffix(fnse.Sel.Name, muStr) {
				if se, ok := fnse.X.(*ast.SelectorExpr); ok {
					if sel, ok := pass.TypesInfo.Selections[se]; ok && sel.Obj() == recvMu {
						return true
					}
				}
			}
		}
	}
	return false
}

// isMutexType is a helper for determining if the object type is a mutex
func isMutexType(t types.Type) bool {
	// The type is a mutex type if it has a `Mutex` suffix
	return strings.HasSuffix(t.String(), "Mutex")
}

// isStaticField returns true if the field can be treated as static and doesn't
// need to be managed under a mutex
func isStaticField(name string) bool {
	return firstWordIs(name, "static") || firstWordIs(name, "atomic")
}

// isSyncObject is a helper to determing if the object is a sync package object
//
// We check for golang's sync packages as well as NebulousLabs TryMutex and
// ThreadGroup packages
func isSyncObject(t types.Type) bool {
	switch t.String() {
	case "sync.Mutex",
		"sync.RWMutex",
		"sync.WaitGroup",
		"gitlab.com/NebulousLabs/Sia/sync.TryMutex",
		"gitlab.com/NebulousLabs/threadgroup.ThreadGroup":
		return true
	}
	return false
}

// managesOwnLocking returns whether a method manages its own locking.
//
// atomic methods use sync.Atomic to be thread safe and don't require a mutex
// extern methods are handled by a mutex external to the struct's primary mutex
// managed, call, and threaded all handle the structs mutex
//   - managed is a synchronous method within a subsystem
//   - call is a synchronous method used by other subsystems
//   - threaded is an asynchronous method
func managesOwnLocking(name string) bool {
	return isManagedExported(name) ||
		firstWordIs(name, "extern") ||
		firstWordIs(name, "managed") ||
		firstWordIs(name, "threaded") ||
		firstWordIs(name, "call")
}

// run implements the analysis interface
func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fd := n.(*ast.FuncDecl)
		if fd.Recv == nil {
			return // not a method
		}
		// if there's no receiver list (e.g. func (Foo) bar()), skip this method, since it obviously
		// can't access a mutex or any other methods
		if len(fd.Recv.List) == 0 || len(fd.Recv.List[0].Names) == 0 {
			return
		}
		recv := pass.TypesInfo.Defs[fd.Recv.List[0].Names[0]]
		mu, ok := containsMutex(pass, recv)
		if !ok {
			return
		}
		checkLockSafety(pass, fd, recv, mu)
	})
	return nil, nil
}
