package main

import (
	"C"
	"log"
	"syscall"
	"unsafe"
)

func attach(pid int) {
	if err := syscall.PtraceAttach(pid); err != nil {
		log.Fatal(err)
	}

	var ws syscall.WaitStatus
	var ru syscall.Rusage

	_, err := syscall.Wait4(pid, &ws, syscall.WSTOPPED, &ru)
	if err != nil {
		log.Fatal(err)
	}

	if err := syscall.PtraceSyscall(pid, 0); err != nil {
		log.Fatal(err)
	}

	var regs syscall.PtraceRegs

	for {
		_, err := syscall.Wait4(pid, &ws, syscall.WSTOPPED, &ru)
		if err != nil {
			log.Fatal(err)
		}
		if err := syscall.PtraceGetRegs(pid, &regs); err != nil {
			log.Fatal(err)
		}
		log.Println(regs.Orig_rax, regs.Rbx, regs.Rcx, regs.Rdx)
		if regs.Orig_rax == 3 {
			p := unsafe.Pointer(&regs.Rcx)
			//log.Println(C.int(regs.Rdx))
			content := C.GoBytes(p, C.int(int(regs.Rdx)))
			log.Println(len(content))
			log.Println(string(content[:200]))

		}
		if err := syscall.PtraceSyscall(pid, 0); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	attach(11508)
}
