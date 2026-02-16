//go:build windows && amd64

package iec61850

// #cgo CFLAGS: -Wno-builtin-declaration-mismatch -I./libiec61850/inc/hal/inc -I./libiec61850/inc/common/inc -I./libiec61850/inc/goose -I./libiec61850/inc/sampled_values -I./libiec61850/inc/iec61850/inc -I./libiec61850/inc/logging -I./libiec61850/inc/mms -I./libiec61850/inc/r_session
// #cgo LDFLAGS: -static -static-libgcc -static-libstdc++ -L${SRCDIR}/libiec61850/lib/win64 ${SRCDIR}/libiec61850/lib/win64/libiec61850.a ${SRCDIR}/libiec61850/lib/win64/libhal.a -lmbedtls -lmbedx509 -lmbedcrypto -lwpcap -lpacket -lws2_32 -liphlpapi -lbcrypt -lmingwex
import "C"
