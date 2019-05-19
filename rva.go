package pe

// An RVA resolver maps a VirtualAddress to a file physical
// address. When the physical file is mapped into memory, sections in
// the file are mapped at different memory addresses. Internally the
// PE file contains pointers to those virtual addresses. This means we
// need to convert these pointers to mapped memory back into the file
// so we can read their data. The RVAResolver is responsible for this
// - it is populated from the header's sections.
type Run struct {
	VirtualAddress  uint32
	VirtualEnd      uint32
	PhysicalAddress uint32
}

type RVAResolver struct {
	// For now very simple O(n) search.
	Runs []*Run
}

func (self *RVAResolver) GetFileAddress(offset uint32) uint32 {
	for _, run := range self.Runs {
		if offset >= run.VirtualAddress &&
			offset < run.VirtualEnd {
			return offset - run.VirtualAddress + run.PhysicalAddress
		}
	}

	return 0
}

func NewRVAResolver(header *IMAGE_NT_HEADERS) *RVAResolver {
	runs := []*Run{}
	for _, section := range header.Sections() {
		if section.SizeOfRawData() == 0 {
			continue
		}

		runs = append(runs, &Run{
			VirtualAddress:  section.VirtualAddress(),
			VirtualEnd:      section.VirtualAddress() + section.SizeOfRawData(),
			PhysicalAddress: section.PointerToRawData(),
		})
	}

	return &RVAResolver{runs}
}
