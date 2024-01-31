package test

import "fmt"

// Origin describes the "source" of a [Test].
type Origin interface {
	AcceptVisitor(OriginVisitor)
	Path() string
	String() string
}

// OriginVisitor is an interface for dispatching based on the concrete type of
// an [Origin].
type OriginVisitor interface {
	VisitDirectory(Directory)
	VisitDocument(Document)
	VisitDocumentLocation(DocumentLocation)
}

// Directory is an [Origin] that refers to a filesystem directory.
type Directory struct {
	Name string
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o Directory) AcceptVisitor(v OriginVisitor) {
	v.VisitDirectory(o)
}

// Path returns the path to the directory.
func (o Directory) Path() string {
	return string(o.Name)
}

func (o Directory) String() string {
	return string(o.Name + "/")
}

// Document is an [Origin] that refers to a specific Markdown document file.
type Document struct {
	Name string
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o Document) AcceptVisitor(v OriginVisitor) {
	v.VisitDocument(o)
}

// Path returns the path to the document file.
func (o Document) Path() string {
	return string(o.Name)
}

func (o Document) String() string {
	return string(o.Name)
}

// DocumentLocation is an [Origin] that describes a specific line within a
// Markdown document.
type DocumentLocation struct {
	File string
	Line int
}

// AcceptVisitor dispatches to the appropriate method on v based on the concrete
// type of o.
func (o DocumentLocation) AcceptVisitor(v OriginVisitor) {
	v.VisitDocumentLocation(o)
}

// Path returns the path to the file.
func (o DocumentLocation) Path() string {
	return o.File
}

func (o DocumentLocation) String() string {
	return fmt.Sprintf("%s:%d", o.File, o.Line)
}
