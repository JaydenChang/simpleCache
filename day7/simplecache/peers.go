package simplecache

import pb "simplecache/simplecache/cachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// Get(group string, k
	Get(in *pb.Request, out *pb.Response) error
}
