package hid

import (
    "github.com/pkg/errors"
    "github.com/sstallion/go-hid"
)

var (
    vendorIDs  = [...]uint16{1356}
    productIDs = [...]uint16{2508, 1476}
)

type Device struct {
    hidDeviceInfo *hid.DeviceInfo
    hidDevice     *hid.Device
}

func (d *Device) VendorID() uint16 {
    return d.hidDeviceInfo.VendorID
}

func (d *Device) ProductID() uint16 {
    return d.hidDeviceInfo.ProductID
}

func (d *Device) Path() string {
    return d.hidDeviceInfo.Path
}

func (d *Device) Release() uint16 {
    return d.hidDeviceInfo.ReleaseNbr
}

func (d *Device) Serial() string {
    return d.hidDeviceInfo.SerialNbr
}

func (d *Device) Manufacturer() string {
    return d.hidDeviceInfo.MfrStr
}

func (d *Device) Product() string {
    return d.hidDeviceInfo.ProductStr
}

func (d *Device) Open() error {
    hidDevice, err := hid.Open(d.VendorID(), d.ProductID(), d.Serial())
    if err != nil {
        return err
    }

    d.hidDevice = hidDevice

    return nil
}

func (d *Device) Close() error {
    return d.hidDevice.Close()
}

func (d *Device) Read(b []byte) (int, error) {
    return d.hidDevice.Read(b)
}

func (d *Device) Write(b []byte) (int, error) {
    return d.hidDevice.Write(b)
}

func (d *Device) GetFeatureReport(code byte) ([]byte, error) {
    var bytes []byte

    switch code {
    case 0x04:
        bytes = make([]byte, 67)
        bytes[0] = code
        _, err := d.hidDevice.GetFeatureReport(bytes)
        if err != nil {
            return nil, err
        }
    default:
        return nil, errors.Errorf("hid: unsupported report code: %v", code)
    }

    return bytes, nil
}

func Find() ([]*Device, error) {
    var devices []*Device

    for _, vendorID := range vendorIDs {
        for _, productID := range productIDs {
            err := hid.Enumerate(vendorID, productID, func(info *hid.DeviceInfo) error {
                devices = append(devices, &Device{hidDeviceInfo: info})
                return nil
            })
            if err != nil {
                return nil, err
            }
        }
    }

    return devices, nil
}
