package test

import (
	"fmt"
)

// Origin describes the "source" of a [Test].
type Origin interface {
	AcceptVisitor(OriginVisitor)
	Path() string
	String() string
}

// OriginVisitor is an interface for dispatching based on the concrete type of
// an [Origin].
type OriginVisitor interface {
	VisitDirectoryOrigin(DirectoryOrigin)
	VisitDocumentOrigin(FileOrigin)
	VisitDocumentLineOrigin(FileLineOrigin)
}

// DirectoryOrigin is an [Origin] that refers to a filesystem directory.
type DirectoryOrigin struct {
	DirPath string
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o DirectoryOrigin) AcceptVisitor(v OriginVisitor) {
	v.VisitDirectoryOrigin(o)
}

// Path returns the path to the directory.
func (o DirectoryOrigin) Path() string {
	return string(o.DirPath)
}

func (o DirectoryOrigin) String() string {
	return string(o.DirPath + "/")
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o DirectoryOrigin) DapperString() string {
	return o.String()
}

// FileOrigin is an [Origin] that refers to a specific file.
type FileOrigin struct {
	FilePath string
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o FileOrigin) AcceptVisitor(v OriginVisitor) {
	v.VisitDocumentOrigin(o)
}

// Path returns the path to the document file.
func (o FileOrigin) Path() string {
	return string(o.FilePath)
}

func (o FileOrigin) String() string {
	return string(o.FilePath)
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o FileOrigin) DapperString() string {
	return o.String()
}

// FileLineOrigin is an [Origin] that describes a specific line within a file.
type FileLineOrigin struct {
	FilePath string
	Line     int
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o FileLineOrigin) AcceptVisitor(v OriginVisitor) {
	v.VisitDocumentLineOrigin(o)
}

// Path returns the path to the file.
func (o FileLineOrigin) Path() string {
	return o.FilePath
}

func (o FileLineOrigin) String() string {
	return fmt.Sprintf("%s:%d", o.FilePath, o.Line)
}

// DapperString returns the string used to represent o in
// [github.com/dogmatiq/dapper] output.
func (o FileLineOrigin) DapperString() string {
	return o.String()
}
