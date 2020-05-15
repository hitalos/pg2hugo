package filesystem

import (
	"context"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"

	"github.com/hitalos/pg2hugo/models"
)

type index struct {
	inode   uint64
	content *models.Content
}

func (i index) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = i.inode
	a.Size = i.content.Size()
	a.Blocks = a.Size / 512
	a.Mode = 0444
	a.Mtime = i.content.LastMod
	a.Ctime = i.content.PublishDate
	return nil
}

func (i index) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	return i, nil
}

func (i index) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	fuseutil.HandleRead(req, resp, []byte(i.content.String()))
	return nil
}
