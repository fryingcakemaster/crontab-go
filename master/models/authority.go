package models

type Authority struct {
	ReadAble  bool
	WriteAble bool
}

func (a *Authority) ChmodAuth(ifRead bool, ifWrite bool)  {
	a.ReadAble = ifRead
	a.WriteAble = ifWrite
}