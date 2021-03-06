/*

Package gogame is a set of functions and modules designed for writing games. Gogame uses SDL2 internally.

This software is free. It's released under the LGPL license. You can create open source and commercial games with it. See the license for full details.

OPENGL is required. Developer libraries of SDL2, SDL2-image and SDL2-TTF are required also.

*/
package gogame

/*
#cgo pkg-config: sdl2 x11
#include "SDL.h"

extern int width;
extern int height;
extern int initWidthAndHeight();
extern int initSDL();
extern SDL_Window * newScreen(char *title, int h, int v);
extern SDL_Renderer * newRenderer( SDL_Window * screen );
extern void setScaleQuality(int n);
extern int isNull(void *pointer);

*/
import "C"
import (
	"errors"
	"unsafe"
)

var screen *C.SDL_Window
var renderer *C.SDL_Renderer

type Color struct {
	R, G, B, A uint8
}

var COLOR_WHITE = &Color{255, 255, 255, 255}
var COLOR_BLACK = &Color{0, 0, 0, 255}
var COLOR_RED = &Color{255, 0, 0, 255}
var COLOR_BLUE = &Color{0, 0, 255, 255}

// Use this function to create a window and a renderer (not visible to user)
func InitXScreenSaver(title string) error {
	if i := C.initWidthAndHeight(); i != 0 {
		return errors.New("Error initializing from X Screensaver window")
	}
	h, v := GetWindowSize()
	return Init(title, h, v)
}

func Init(title string, h, v int) error {
	if i := C.initSDL(); i != 0 {
		return errors.New("Error initializing SDL")
	}
	screen = C.newScreen(C.CString(title), C.int(h), C.int(v))
	if C.isNull(unsafe.Pointer(screen)) == 1 {
		return errors.New("Error initalizing SCREEN")
	}
	renderer = C.newRenderer(screen)
	if C.isNull(unsafe.Pointer(renderer)) == 1 {
		return errors.New("Error initalizing RENDERER")
	}
	if screen == nil || renderer == nil {
		return errors.New("Error on initializing SDL2")
	}
	return nil
}

func SetScaleQuality(n int) {
	C.setScaleQuality(C.int(n))
}

// Full Screen mode
func SetFullScreen(fs bool) {
	if fs {
		C.SDL_SetWindowFullscreen(screen, C.SDL_WINDOW_FULLSCREEN_DESKTOP)
	} else {
		C.SDL_SetWindowFullscreen(screen, 0)
	}
}

// Get window size
func GetWindowSize() (int, int) {
	w, h := int(C.width), int(C.height)
	if w == 0 || h == 0 {
		var w, h C.int
		C.SDL_GetWindowSize(screen, &w, &h)
		return int(w), int(h)
	}
	return w, h
}

// Set window size
func SetWindowSize(h, v int) {
	C.SDL_SetWindowSize(screen, C.int(h), C.int(v))
}

// Set a device independent resolution for rendering
func SetLogicalSize(h, v int) {
	C.SDL_RenderSetLogicalSize(renderer, C.int(h), C.int(v))
}

// Destroys renderer and window
func Quit() {
	C.SDL_DestroyRenderer(renderer)
	C.SDL_DestroyWindow(screen)
	C.SDL_Quit()
}

// Clear the current rendering target with black color
func RenderClear() {
	C.SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0)
	C.SDL_RenderClear(renderer)
}

// Update the screen with rendering performed
func RenderPresent() {
	C.SDL_RenderPresent(renderer)
}

// Wait specified number of milliseconds before returning
func Delay(s int) {
	C.SDL_Delay(C.Uint32(s))
}

// Draw pixel at position x,y
func DrawPixel(x, y int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	C.SDL_RenderDrawPoint(renderer, C.int(x), C.int(y))
}

// Draw line
func DrawLine(x1, y1, x2, y2 int, color *Color) {
	C.SDL_SetRenderDrawColor(renderer, C.Uint8(color.R), C.Uint8(color.G), C.Uint8(color.B), C.Uint8(color.A))
	C.SDL_RenderDrawLine(renderer, C.int(x1), C.int(y1), C.int(x2), C.int(y2))
}
