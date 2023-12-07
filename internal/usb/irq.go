// https://github.com/usbarmory/GoKey
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

//go:build tamago && arm
// +build tamago,arm

package usb

import (
	"log"

	"github.com/usbarmory/tamago/arm"
	"github.com/usbarmory/tamago/soc/nxp/imx6ul"
	"github.com/usbarmory/tamago/soc/nxp/usb"
)

func StartInterruptHandler(port *usb.USB) {
	if port == nil {
		return
	}

	imx6ul.GIC.Init(true, false)
	imx6ul.GIC.EnableInterrupt(port.IRQ, true)

	port.EnableInterrupt(usb.IRQ_URI) // reset
	port.EnableInterrupt(usb.IRQ_PCI) // port change detect
	port.EnableInterrupt(usb.IRQ_UI)  // transfer completion

	arm.RegisterInterruptHandler()

	for {
		arm.WaitInterrupt()

		imx6ul.SetARMFreq(imx6ul.FreqMax)
		defer imx6ul.SetARMFreq(imx6ul.FreqLow)

		irq, end := imx6ul.GIC.GetInterrupt(true)

		if end != nil {
			close(end)
		}

		switch {
		case irq == port.IRQ:
			port.ServiceInterrupts()
		default:
			log.Printf("internal error, unexpected IRQ %d", irq)
		}
	}
}
