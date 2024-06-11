package main 

type ErrMsg struct{
  isFatal bool
  err error
}

func (e ErrMsg) Error() string {
  return e.err.Error();
}

func (e ErrMsg) IsFatal() bool {
  return e.isFatal;
}

