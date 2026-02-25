package protojson

import "google.golang.org/protobuf/reflect/protoreflect"

// arrayUnpopulatedFieldRanger wraps a protoreflect.Message and modifies its Range
// method to additionally iterate over unpopulated repeated fields, emitting them
// as empty arrays. All other unpopulated fields are skipped (default behavior).
type arrayUnpopulatedFieldRanger struct {
	protoreflect.Message
}

func (m arrayUnpopulatedFieldRanger) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		if m.Has(fd) || fd.ContainingOneof() != nil {
			continue // ignore populated fields and fields within a oneof
		}
		if fd.Cardinality() != protoreflect.Repeated {
			continue // skip all non-repeated unpopulated fields
		}
		if !f(fd, m.Get(fd)) { // m.Get on an unset repeated field returns an empty list
			return
		}
	}
	m.Message.Range(f)
}
