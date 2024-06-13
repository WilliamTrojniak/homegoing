package main 

type ErrMsg struct{
  isFatal bool
  error
}

func (e ErrMsg) IsFatal() bool {
  return e.isFatal;
}

