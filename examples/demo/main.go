package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/kpeu3i/gods4"
    "github.com/kpeu3i/gods4/led"
)

func main() {
    // Find all controllers connected to your machine via USB or Bluetooth
    controllers := gods4.Find()
    if len(controllers) == 0 {
        panic("No connected DS4 controllers found")
    }

    // Select first controller from the list
    controller := controllers[0]

    // Connect to the controller
    err := controller.Connect()
    if err != nil {
        panic(err)
    }

    log.Printf("* Controller #1 | %-10s | name: %s, connection: %s\n", "Connect", controller, controller.ConnectionType())

    // Disconnect controller when a program is terminated
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-signals
        err := controller.Disconnect()
        if err != nil {
            panic(err)
        }
        log.Printf("* Controller #1 | %-10s | bye!\n", "Disconnect")
    }()

    // Register callback for "BatteryUpdate" event
    controller.On(gods4.EventBatteryUpdate, func(data interface{}) error {
        battery := data.(gods4.Battery)
        log.Printf("* Controller #1 | %-10s | capacity: %v%%, charging: %v, cable: %v\n",
            "Battery",
            battery.Capacity,
            battery.IsCharging,
            battery.IsCableConnected,
        )

        return nil
    })

    // Register callback for "CrossPress" event
    controller.On(gods4.EventCrossPress, func(data interface{}) error {
        log.Printf("* Controller #1 | %-10s | state: press\n", "Cross")

        return nil
    })

    // Register callback for "CrossRelease" event
    controller.On(gods4.EventCrossRelease, func(data interface{}) error {
        log.Printf("* Controller #1 | %-10s | state: release\n", "Cross")

        return nil
    })

    // Register callback for "RightStickMove" event
    controller.On(gods4.EventRightStickMove, func(data interface{}) error {
        stick := data.(gods4.Stick)
        log.Printf("* Controller #1 | %-10s | x: %v, y: %v\n", "RightStick", stick.X, stick.Y)

        return nil
    })

    // Enable LED (yellow) with flash
    err = controller.Led(led.Yellow())
    if err != nil {
        panic(err)
    }

    // Start listening for controller events
    err = controller.Listen()
    if err != nil {
        panic(err)
    }
}
