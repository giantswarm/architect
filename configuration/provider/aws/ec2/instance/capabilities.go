// Package instance provides configuration structures for the EC2 instance
// specific settings.
package instance

// capabilities is a private type to ensure only instance capabilities defined
// in this package can be applied to installation configurations. That prevents
// other packages screwing around with instance capability configurations.
type capabilities struct {
	CPUCores      int     `json:"cpu_cores"`
	MemorySizeGB  float64 `json:"memory_size_gb"`
	StorageSizeGB float64 `json:"storage_size_gb"`
}

// Capabilities returns the full list of capabilities of all available instance
// types as defined by the constants of this package. Using this method ensures
// that always the same immutable list will be used when referencing
// capabilities of all available instance types EC2 provides. Note that this
// method was generated using the following script.
//
//    https://gist.github.com/xh3b4sd/55b0bff0c4faf0a78a77b1361c393052
//
func Capabilities() map[kind]capabilities {
	return map[kind]capabilities{
		TypeC1Medium: {
			CPUCores:      2,
			MemorySizeGB:  1.7,
			StorageSizeGB: 350,
		},
		TypeC1XLarge: {
			CPUCores:      8,
			MemorySizeGB:  7,
			StorageSizeGB: 420,
		},
		TypeC32XLarge: {
			CPUCores:      8,
			MemorySizeGB:  15,
			StorageSizeGB: 80,
		},
		TypeC34XLarge: {
			CPUCores:      16,
			MemorySizeGB:  30,
			StorageSizeGB: 160,
		},
		TypeC38XLarge: {
			CPUCores:      32,
			MemorySizeGB:  60,
			StorageSizeGB: 320,
		},
		TypeC3Large: {
			CPUCores:      2,
			MemorySizeGB:  3.75,
			StorageSizeGB: 16,
		},
		TypeC3XLarge: {
			CPUCores:      4,
			MemorySizeGB:  7.5,
			StorageSizeGB: 40,
		},
		TypeC42XLarge: {
			CPUCores:      8,
			MemorySizeGB:  15,
			StorageSizeGB: 0,
		},
		TypeC44XLarge: {
			CPUCores:      16,
			MemorySizeGB:  30,
			StorageSizeGB: 0,
		},
		TypeC48XLarge: {
			CPUCores:      36,
			MemorySizeGB:  60,
			StorageSizeGB: 0,
		},
		TypeC4Large: {
			CPUCores:      2,
			MemorySizeGB:  3.75,
			StorageSizeGB: 0,
		},
		TypeC4XLarge: {
			CPUCores:      4,
			MemorySizeGB:  7.5,
			StorageSizeGB: 0,
		},
		TypeCC28XLarge: {
			CPUCores:      32,
			MemorySizeGB:  60.5,
			StorageSizeGB: 840,
		},
		TypeCG14XLarge: {
			CPUCores:      16,
			MemorySizeGB:  22.5,
			StorageSizeGB: 840,
		},
		TypeCR18XLarge: {
			CPUCores:      32,
			MemorySizeGB:  244,
			StorageSizeGB: 120,
		},
		TypeD22XLarge: {
			CPUCores:      8,
			MemorySizeGB:  61,
			StorageSizeGB: 2000,
		},
		TypeD24XLarge: {
			CPUCores:      16,
			MemorySizeGB:  122,
			StorageSizeGB: 2000,
		},
		TypeD28XLarge: {
			CPUCores:      36,
			MemorySizeGB:  244,
			StorageSizeGB: 2000,
		},
		TypeD2XLarge: {
			CPUCores:      4,
			MemorySizeGB:  30.5,
			StorageSizeGB: 2000,
		},
		TypeF116XLarge: {
			CPUCores:      64,
			MemorySizeGB:  976,
			StorageSizeGB: 940,
		},
		TypeF12XLarge: {
			CPUCores:      8,
			MemorySizeGB:  122,
			StorageSizeGB: 470,
		},
		TypeG22XLarge: {
			CPUCores:      8,
			MemorySizeGB:  15,
			StorageSizeGB: 60,
		},
		TypeG28XLarge: {
			CPUCores:      32,
			MemorySizeGB:  60,
			StorageSizeGB: 120,
		},
		TypeHI14XLarge: {
			CPUCores:      16,
			MemorySizeGB:  60.5,
			StorageSizeGB: 1024,
		},
		TypeHS18XLarge: {
			CPUCores:      16,
			MemorySizeGB:  117,
			StorageSizeGB: 2000,
		},
		TypeI22XLarge: {
			CPUCores:      8,
			MemorySizeGB:  61,
			StorageSizeGB: 800,
		},
		TypeI24XLarge: {
			CPUCores:      16,
			MemorySizeGB:  122,
			StorageSizeGB: 800,
		},
		TypeI28XLarge: {
			CPUCores:      32,
			MemorySizeGB:  244,
			StorageSizeGB: 800,
		},
		TypeI2XLarge: {
			CPUCores:      4,
			MemorySizeGB:  30.5,
			StorageSizeGB: 800,
		},
		TypeI316XLarge: {
			CPUCores:      64,
			MemorySizeGB:  488,
			StorageSizeGB: 1900,
		},
		TypeI32XLarge: {
			CPUCores:      8,
			MemorySizeGB:  61,
			StorageSizeGB: 1900,
		},
		TypeI34XLarge: {
			CPUCores:      16,
			MemorySizeGB:  122,
			StorageSizeGB: 1900,
		},
		TypeI38XLarge: {
			CPUCores:      32,
			MemorySizeGB:  244,
			StorageSizeGB: 1900,
		},
		TypeI3Large: {
			CPUCores:      2,
			MemorySizeGB:  15.25,
			StorageSizeGB: 475,
		},
		TypeI3XLarge: {
			CPUCores:      4,
			MemorySizeGB:  30.5,
			StorageSizeGB: 950,
		},
		TypeM1Large: {
			CPUCores:      2,
			MemorySizeGB:  7.5,
			StorageSizeGB: 420,
		},
		TypeM1Medium: {
			CPUCores:      1,
			MemorySizeGB:  3.75,
			StorageSizeGB: 410,
		},
		TypeM1Small: {
			CPUCores:      1,
			MemorySizeGB:  1.7,
			StorageSizeGB: 160,
		},
		TypeM1XLarge: {
			CPUCores:      4,
			MemorySizeGB:  15,
			StorageSizeGB: 420,
		},
		TypeM22XLarge: {
			CPUCores:      4,
			MemorySizeGB:  34.2,
			StorageSizeGB: 850,
		},
		TypeM24XLarge: {
			CPUCores:      8,
			MemorySizeGB:  68.4,
			StorageSizeGB: 840,
		},
		TypeM2XLarge: {
			CPUCores:      2,
			MemorySizeGB:  17.1,
			StorageSizeGB: 420,
		},
		TypeM32XLarge: {
			CPUCores:      8,
			MemorySizeGB:  30,
			StorageSizeGB: 80,
		},
		TypeM3Large: {
			CPUCores:      2,
			MemorySizeGB:  7.5,
			StorageSizeGB: 32,
		},
		TypeM3Medium: {
			CPUCores:      1,
			MemorySizeGB:  3.75,
			StorageSizeGB: 4,
		},
		TypeM3XLarge: {
			CPUCores:      4,
			MemorySizeGB:  15,
			StorageSizeGB: 40,
		},
		TypeM410XLarge: {
			CPUCores:      40,
			MemorySizeGB:  160,
			StorageSizeGB: 0,
		},
		TypeM416XLarge: {
			CPUCores:      64,
			MemorySizeGB:  256,
			StorageSizeGB: 0,
		},
		TypeM42XLarge: {
			CPUCores:      8,
			MemorySizeGB:  32,
			StorageSizeGB: 0,
		},
		TypeM44XLarge: {
			CPUCores:      16,
			MemorySizeGB:  64,
			StorageSizeGB: 0,
		},
		TypeM4Large: {
			CPUCores:      2,
			MemorySizeGB:  8,
			StorageSizeGB: 0,
		},
		TypeM4XLarge: {
			CPUCores:      4,
			MemorySizeGB:  16,
			StorageSizeGB: 0,
		},
		TypeP216XLarge: {
			CPUCores:      64,
			MemorySizeGB:  732,
			StorageSizeGB: 0,
		},
		TypeP28XLarge: {
			CPUCores:      32,
			MemorySizeGB:  488,
			StorageSizeGB: 0,
		},
		TypeP2XLarge: {
			CPUCores:      4,
			MemorySizeGB:  61,
			StorageSizeGB: 0,
		},
		TypeR32XLarge: {
			CPUCores:      8,
			MemorySizeGB:  61,
			StorageSizeGB: 160,
		},
		TypeR34XLarge: {
			CPUCores:      16,
			MemorySizeGB:  122,
			StorageSizeGB: 320,
		},
		TypeR38XLarge: {
			CPUCores:      32,
			MemorySizeGB:  244,
			StorageSizeGB: 320,
		},
		TypeR3Large: {
			CPUCores:      2,
			MemorySizeGB:  15.25,
			StorageSizeGB: 32,
		},
		TypeR3XLarge: {
			CPUCores:      4,
			MemorySizeGB:  30.5,
			StorageSizeGB: 80,
		},
		TypeR416XLarge: {
			CPUCores:      64,
			MemorySizeGB:  488,
			StorageSizeGB: 0,
		},
		TypeR42XLarge: {
			CPUCores:      8,
			MemorySizeGB:  61,
			StorageSizeGB: 0,
		},
		TypeR44XLarge: {
			CPUCores:      16,
			MemorySizeGB:  122,
			StorageSizeGB: 0,
		},
		TypeR48XLarge: {
			CPUCores:      32,
			MemorySizeGB:  244,
			StorageSizeGB: 0,
		},
		TypeR4Large: {
			CPUCores:      2,
			MemorySizeGB:  15.25,
			StorageSizeGB: 0,
		},
		TypeR4XLarge: {
			CPUCores:      4,
			MemorySizeGB:  30.5,
			StorageSizeGB: 0,
		},
		TypeT1Micro: {
			CPUCores:      1,
			MemorySizeGB:  0.613,
			StorageSizeGB: 0,
		},
		TypeT22XLarge: {
			CPUCores:      8,
			MemorySizeGB:  32,
			StorageSizeGB: 0,
		},
		TypeT2Large: {
			CPUCores:      2,
			MemorySizeGB:  8,
			StorageSizeGB: 0,
		},
		TypeT2Medium: {
			CPUCores:      2,
			MemorySizeGB:  4,
			StorageSizeGB: 0,
		},
		TypeT2Micro: {
			CPUCores:      1,
			MemorySizeGB:  1,
			StorageSizeGB: 0,
		},
		TypeT2Nano: {
			CPUCores:      1,
			MemorySizeGB:  0.5,
			StorageSizeGB: 0,
		},
		TypeT2Small: {
			CPUCores:      1,
			MemorySizeGB:  2,
			StorageSizeGB: 0,
		},
		TypeT2XLarge: {
			CPUCores:      4,
			MemorySizeGB:  16,
			StorageSizeGB: 0,
		},
		TypeX116XLarge: {
			CPUCores:      64,
			MemorySizeGB:  976,
			StorageSizeGB: 1920,
		},
		TypeX132XLarge: {
			CPUCores:      128,
			MemorySizeGB:  1952,
			StorageSizeGB: 1920,
		},
	}
}
