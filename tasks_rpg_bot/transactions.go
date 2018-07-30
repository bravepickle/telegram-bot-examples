package main

type Transactional interface {
	GetName() string
	//// GetData returns transactional data
	GetData() map[string]interface{}
	SetData(name string, value interface{})
	Init()
	Current() int
	//GetSteps() []string
	Next()
	Prev()
	Reset()
	Run() bool
	//Complete() bool
	//Commit() bool
}

type TransactionalStep interface {
	GetName() string
	Run(t Transactional) bool
	Revert(t Transactional) // revert run step
}

// =========== TransactionStruct
type TransactionStruct struct {
	currentStepIndex int
	//steps            map[int]TransactionalStep
	steps []TransactionalStep
	data  map[string]interface{}
}

func (t *TransactionStruct) Next() {
	t.currentStepIndex += 1
}

func (t *TransactionStruct) Prev() {
	t.currentStepIndex -= 1
}

func (t *TransactionStruct) Reset() {
	t.currentStepIndex = 0
	t.data = make(map[string]interface{})
	//t.steps = make(map[int]TransactionalStep)
	//t.steps = make(map[int]TransactionalStep)
}

func (t *TransactionStruct) Current() int {
	return t.currentStepIndex
}

func (t *TransactionStruct) Init() {
}

func (t *TransactionStruct) Run() bool {
	//t.RunStep()

	return true
}

func (t *TransactionStruct) RunStep() {
	//if step, ok := t.steps[t.currentStepIndex]; !ok {
	//if step := t.steps[t.currentStepIndex]; !ok {
	if step := t.steps[t.currentStepIndex]; step != nil {
		logger.Error(`No steps initialized for transaction %s`, t.GetName())
	} else {
		step.Run(t)
	}
}

func (t *TransactionStruct) GetName() string {
	return `undefined` // override in children
}

func (t *TransactionStruct) GetData() map[string]interface{} {
	return t.data
}

func (t *TransactionStruct) SetData(name string, value interface{}) {
	t.data[name] = value
}

// =========== AddTaskTransactionStruct

type AddTaskTransactionStruct struct {
	TransactionStruct
}

func (t *AddTaskTransactionStruct) GetName() string {
	return `add-task`
}

func (t *AddTaskTransactionStruct) Init() {
	t.Reset()

	t.steps = append(t.steps, SetTitleStep{})
}

func (t *AddTaskTransactionStruct) Run() bool {
	t.RunStep()

	return true
}

// =========== SetTitleStep
type SetTitleStep struct {
	TransactionalStep
}

//GetName() string
//Run(t Transactional) bool
//Revert() // revert runned step
func (t SetTitleStep) GetName() string {
	return `set-title`
}
func (t SetTitleStep) Run(tr Transactional) bool {
	// TODO: generate message here to input
	//tr.SetData(`title`, `abcddsd // `)

	return true
}
func (t SetTitleStep) Revert(tr Transactional) {
	tr.SetData(`title`, nil) // TODO: revert properly
}

func NewAddTaskTransaction() *AddTaskTransactionStruct {
	var transaction AddTaskTransactionStruct
	transaction.Init()

	return &transaction
}
