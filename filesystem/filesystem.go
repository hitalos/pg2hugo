package filesystem

import (
	"context"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/hitalos/pg2hugo/models"
)

var inodeCounter uint64 = 1

// FS implements filesystem
type FS struct {
	dirs map[string]*dir
}

// NewFS create a new FS instance
func NewFS() (*FS, error) {
	if err := models.Connect(); err != nil {
		return nil, err
	}
	fs := FS{}
	if err := readAll(&fs); err != nil {
		return nil, err
	}
	log.Println("Filesystem mounted and ready to use")
	return &fs, nil
}

func readAll(f *FS) error {
	contents, err := models.ReadAllContents()
	if err != nil {
		return err
	}
	inodeCounter = 1
	f.dirs = make(map[string]*dir, len(contents))
	for _, content := range contents {
		inodeCounter++
		entry := fuse.Dirent{
			Inode: inodeCounter,
			Type:  fuse.DT_Dir,
			Name:  content.Path,
		}
		inodeCounter += uint64(len(content.Resources)) + 1 // one more for "index.md"

		d := dir{content, entry}
		f.dirs[content.Path] = &d
	}
	return nil
}

// Root implements Root of filesystem
func (f *FS) Root() (fs.Node, error) {
	return f, nil
}

// Attr attributes for root dir
func (f *FS) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
	a.Mode = os.ModeDir | 0o555
	a.Size = 4096
	return nil
}

// ReadDirAll list content of filesystem
func (f *FS) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list := make([]fuse.Dirent, len(f.dirs)+2)
	list[0] = fuse.Dirent{Name: ".", Type: fuse.DT_Dir}
	list[1] = fuse.Dirent{Name: "..", Inode: 1, Type: fuse.DT_Dir}
	i := 2
	for _, dir := range f.dirs {
		list[i] = dir.entry
		i++
	}
	return list, nil
}

// Lookup checks an entry on filesystem
func (f FS) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if c, ok := f.dirs[name]; ok {
		return c, nil
	}
	return nil, syscall.ENOENT
}
