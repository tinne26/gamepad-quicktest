# gamepad-quicktest

Very basic program to test a gamepad with [Ebitengine](https://github.com/hajimehoshi/ebiten) (v2.5.0 or higher required):
```
go run github.com/tinne26/gamepad-quicktest@latest
```

Interactive web version also available at [tinne26.github.io/gamepad-quicktest](https://tinne26.github.io/gamepad-quicktest/).

Sample visual output:

![](https://github.com/tinne26/gamepad-quicktest/blob/main/sample.png?raw=true)

(Notice that the official example on [ebiten/examples/gamepad](https://github.com/hajimehoshi/ebiten/tree/main/examples/gamepad) is also very complete and you can test it [interactively](https://ebitengine.org/en/examples/gamepad.html) too)


## Platform differences

Notice that the behavior for gamepads can vary quite a lot between platforms and different browsers:
- Buttons with an analog response may be detected as binary buttons on some platforms or as axes on others. This is common with the bottom triggers on the front face of the controllers (also known as `R2/L2`, `ZR/ZL`, `RT/LT`... and many other names). Be aware of this, it can cause some problems. Notice also that the threshold for deciding if an analog trigger is pressed or not is fairly arbitrary. Ebitengine uses a threshold of [approximately 0.12](https://github.com/hajimehoshi/ebiten/blob/v2.5.4/internal/gamepaddb/gamepaddb.go#LL526C1-L526C44).
- Vibration support is still fairly limited on Ebitengine, but also on many browsers. Chrome seems to do quite well. If vibration doesn't work for your gamepad on one browser... try another. In my experience, Firefox is not great on this front.
