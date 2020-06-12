package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/hitalos/pg2hugo/filesystem"
)

var preload = flag.Bool("p", false, "Preload all resources")

func usage() {
	fmt.Printf("Use %s [-p] <MOUNTPOINT>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *preload {
		os.Setenv("PRELOAD", "true")
	}
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
	}

	FS, err := filesystem.NewFS()
	if err != nil {
		log.Fatalf("Error creating FS: %s\n", err)
	}

	c, err := fuse.Mount(
		args[0],
		fuse.FSName("Hugo content folder"),
		fuse.Subtype("pg2hugo"),
	)
	if err != nil {
		log.Fatalf("error mounting filesystem: %s", err)
	}
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt)
	go gracefulStop(cc, args[0])

	defer func() {
		if err := c.Close(); err != nil {
			log.Fatalf("Error closing FUSE: %s\n", err)
		}
	}()

	if err := fs.Serve(c, FS); err != nil {
		log.Fatalf("Error serving filesystem: %s\n", err)
	}

	<-c.Ready
	if c.MountError != nil {
		log.Fatalln(c.MountError)
	}
	log.Println("Filesystem umounted correctly")
}

func gracefulStop(cc chan os.Signal, mountpoint string) {
	<-cc
	log.Printf("unmounting %q\n", mountpoint)
	err := fuse.Unmount(mountpoint)
	if err != nil {
		log.Fatalf("error unmounting")
	}
}
