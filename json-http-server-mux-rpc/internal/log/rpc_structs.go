package log

type AppendArgs struct {
	Record Record
}

type AppendReply struct {
	Offset uint64
}

type FetchArgs struct {
	Offset uint64
}

type FetchReply struct {
	Record Record
}
