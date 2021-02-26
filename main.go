package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func main() {
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	fmt.Printf("Bus: %v ", b.String())
	// Dev is a valid conn.Conn.
	// 0x12 is the Adafruit PMSA003l
	d := &i2c.Dev{Addr: 0x12, Bus: b}
	fmt.Printf("Device: %v ", d.String())

	// Prints out the gpio pin used.
	if p, ok := b.(i2c.Pins); ok {
		fmt.Printf("SDA: %s ", p.SDA())
		fmt.Printf("SCL: %s\n", p.SCL())
	}

	// Send a command 0x10 and expect a 5 bytes reply.
	//write := []byte{0x10}
	for {
		read := make([]byte, 32)
		if err := d.Tx(nil, read); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Raw Packet: %v\n", read)

		fmt.Println("Frame Length: ", binary.BigEndian.Uint16(read[2:4]))

		fmt.Println("PM1.0: ", binary.BigEndian.Uint16(read[4:6]))
		fmt.Println("PM2.5: ", binary.BigEndian.Uint16(read[6:8]))
		fmt.Println("PM10 : ", binary.BigEndian.Uint16(read[8:10]))

		fmt.Println("PM1.0: ", binary.BigEndian.Uint16(read[10:12]))
		fmt.Println("PM2.5: ", binary.BigEndian.Uint16(read[12:14]))
		fmt.Println("PM10 : ", binary.BigEndian.Uint16(read[14:16]))

		fmt.Println("Particles >  0.3 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[16:18]))
		fmt.Println("Particles >  0.5 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[18:20]))
		fmt.Println("Particles >  1.0 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[20:22]))
		fmt.Println("Particles >  2.5 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[22:24]))
		fmt.Println("Particles >  5.0 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[24:26]))
		fmt.Println("Particles > 10.0 µ𝑚 / 0.1L : ", binary.BigEndian.Uint16(read[26:28]))

		fmt.Println("Version : ", int(read[29]))

		fmt.Println("Error : ", int(read[30]))

		var checksum int
		for _, v := range read[0:30] {
			checksum = checksum + int(v)
		}

		if checksum == int(binary.BigEndian.Uint16(read[30:])) {
			fmt.Println("Checksum validated!")
		} else {
			fmt.Println("Checksum did not match!")
			fmt.Println("Checksum (calculated) : ", checksum)
			fmt.Println("Checksum (in packet)  : ", binary.BigEndian.Uint16(read[30:]))
		}

		time.Sleep(3 * time.Second)
	}
}
