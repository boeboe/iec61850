//go:build linux && arm64

package iec61850

// #cgo CFLAGS: -I./libiec61850/inc/hal/inc -I./libiec61850/inc/common/inc -I./libiec61850/inc/goose -I./libiec61850/inc/sampled_values -I./libiec61850/inc/iec61850/inc -I./libiec61850/inc/logging -I./libiec61850/inc/mms -I./libiec61850/inc/r_session
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./libiec61850/lib/linux_armv8 -Wl,--start-group -liec61850 -lhal -lmbedtls -lmbedx509 -lmbedcrypto -Wl,--end-group -lpthread
import "C"
