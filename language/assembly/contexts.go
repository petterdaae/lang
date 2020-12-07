package assembly

import "fmt"

type ContextElementType int

const (
	StackNumElem ContextElementType = iota
	StackBoolElem
	HeapElem
	ProcedureElem
	InvalidElem
)

type Contexts struct {
	stack []*context
	procedureNameCounter int
}

type context struct {
	procedures map[string]string
	Stack      map[string]int
	stackTypes map[string]ContextElementType
	Procedure  *Procedure
}

func NewContexts() *Contexts {
	return &Contexts{
		stack: []*context{
			&context{
				procedures: make(map[string]string),
				Stack:      make(map[string]int),
				stackTypes: make(map[string]ContextElementType),
				Procedure:  nil,
			},
		},
		procedureNameCounter: 0,
	}
}

func (contexts *Contexts) Push(context *context) {
	contexts.stack = append(contexts.stack, context)
}

func (contexts *Contexts) NewContext(newProcedure bool, stackSize int) (*context, []string, int) {
	var operations []string
	initStackSize := stackSize
	current := contexts.Peek()
	stackCopy := make(map[string]int)
	for k, _ := range current.Stack {
		_, address, _ := contexts.Get(k, stackSize)
		operations = append(operations, fmt.Sprintf("mov rax, %s", address))
		operations = append(operations, fmt.Sprintf("push rax"))
		stackSize++
		stackCopy[k] = stackSize
	}
	proceduresCopy := make(map[string]string)
	for k, v := range current.procedures {
		proceduresCopy[k] = v
	}
	stackTypesCopy := make(map[string]ContextElementType)
	for k, v := range current.stackTypes {
		stackTypesCopy[k] = v
	}

	var proc *Procedure
	if newProcedure {
		contexts.procedureNameCounter++
		proc = &Procedure{
			Name:                     fmt.Sprintf("proc%d", contexts.procedureNameCounter),
			Operations:               []string{},
			StackSizeWhenInitialized: stackSize,
		}
		proc.start(stackSize)
	}

	new := &context{
		Stack:      stackCopy,
		procedures: proceduresCopy,
		Procedure:  proc,
		stackTypes: stackTypesCopy,
	}

	return new, operations, stackSize - initStackSize
}

func (contexts *Contexts) Size() int {
	return len(contexts.stack)
}

func (contexts *Contexts) Peek() *context {
	size := contexts.Size()
	if size == 0 {
		return nil
	}
	return contexts.stack[size-1]
}

func (contexts *Contexts) Pop(stackSize int) (int, *context) {
	result := contexts.Peek()
	if result == nil {
		return 0, nil
	}
	size := contexts.Size()
	contexts.stack = contexts.stack[:size-1]
	pops := 0
	if result.Procedure != nil {
		pops += result.Procedure.end(stackSize)
	}
	return pops, result
}

func (contexts *Contexts) getFromContext(context *context, name string, stackSize int) (ContextElementType, string, error) {
	if context == nil {
		return InvalidElem, "", fmt.Errorf("context stack is empty")
	}

	procedure := contexts.GetTopProcedure()

	stack, ok := context.Stack[name]
	kind, _ := context.stackTypes[name]
	if ok {
		diff := (stackSize - stack) * 8

		if procedure != nil {
			if stack > procedure.StackSizeWhenInitialized {
				return kind, fmt.Sprintf("[rsp+%d]", diff), nil
			}

			return kind, fmt.Sprintf("[rsp+rcx+%d+8]", diff), nil
		}

		return kind, fmt.Sprintf("[rsp+%d]", diff), nil
	}

	proc, ok := context.procedures[name]
	if ok {
		return ProcedureElem, proc, nil
	}

	return InvalidElem, "", fmt.Errorf("could not get '%s' from context", name)
}

func (contexts *Contexts) GetTopProcedure() *Procedure {
	for i := contexts.Size() - 1; i >= 0; i-- {
		current := contexts.stack[i].Procedure
		if current != nil {
			return current
		}
	}
	return nil
}

func (contexts *Contexts) Get(name string, stackSize int) (ContextElementType, string, error) {
	context := contexts.Peek()
	return contexts.getFromContext(context, name, stackSize)
}

func (contexts *Contexts) StackInsert(name string, value string, stackSize int, kind ContextElementType) ([]string, int, error) {
	var operations []string
	initStackSize := stackSize
	top := contexts.Peek()
	if top == nil {
		return nil, 0, fmt.Errorf("context stack is empty")
	}

	_, ok := top.Stack[name]
	if !ok {
		operations = append(operations, fmt.Sprintf("push %s", value))
		stackSize++
		top.Stack[name] = stackSize
		top.stackTypes[name] = kind
		return operations, stackSize - initStackSize, nil
	}

	if top.Procedure == nil {
		for i := contexts.Size() - 1; i >= 0; i-- {
			current := contexts.stack[i]
			_, ok := current.Stack[name]
			if !ok {
				break
			}

			_, address, _ := contexts.getFromContext(current, name, stackSize)
			operations = append(operations, fmt.Sprintf("mov %s, %s", address, value))

			if current.Procedure != nil {
				operations = append(operations, fmt.Sprintf("push %s", value))
				stackSize++
				current.Stack[name] = stackSize
				current.stackTypes[name] = kind
				return operations, stackSize - initStackSize, nil
				break
			}
		}
	} else {
		_, address, _ := contexts.getFromContext(top, name, stackSize)
		top.stackTypes[name] = kind
		operations = append(operations, fmt.Sprintf("mov %s, %s", address, value))
		return operations, stackSize - initStackSize, nil
	}

	return operations, initStackSize - stackSize, nil
}

func (contexts *Contexts) ProcInsert(name string, alias string) error {
	top := contexts.Peek()
	if top == nil {
		return fmt.Errorf("context stack is empty")
	}

	_, ok := top.procedures[name]
	if !ok {
		top.procedures[name] = alias
		return nil
	}

	if top.Procedure == nil {
		for i := contexts.Size() - 1; i >= 0; i-- {
			current := contexts.stack[i]
			_, ok := current.procedures[name]
			if !ok {
				break
			}

			current.procedures[name] = alias

			if current.Procedure != nil {
				break
			}
		}
	}

	return nil
}
