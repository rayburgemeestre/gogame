package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

int getEventType(SDL_Event *e) {
    return e->type;
}

int getKeyCode(SDL_Event *e) {
    return e->key.keysym.sym;
}

int isKeyRepeat(SDL_Event *e) {
    return e->key.repeat;
}

int isKeyPressed(int kcode) {
    const Uint8 * state = SDL_GetKeyboardState(NULL);
    int sc = SDL_GetScancodeFromKey(kcode);
    return state[sc];
}

int getWindowEvent(SDL_Event *e) {
    return e->window.event;
}

void getWinData(SDL_Event *e, int *w, int *h) {
    *w = e->window.data1;
    *h = e->window.data2;
}

void getMouseCoords(SDL_Event *e, int *x, int *y, int *down) {
    *x = e->button.x;
    *y = e->button.y;
    if (e->button.state == SDL_PRESSED) {
        *down = 1;
    } else {
        *down = 0;
    }
}

void getMouseWheel(SDL_Event *e, int *x, int *y) {
    *x = e->wheel.x;
    *y = e->wheel.y;
}

*/
import "C"

const (
	K_LEFT   = C.SDLK_LEFT
	K_RIGHT  = C.SDLK_RIGHT
	K_UP     = C.SDLK_UP
	K_DOWN   = C.SDLK_DOWN
	K_SPACE  = C.SDLK_SPACE
	K_ESC    = C.SDLK_ESCAPE
	K_RETURN = C.SDLK_RETURN
	K_P      = C.SDLK_p
	K_I      = C.SDLK_i
	K_M      = C.SDLK_m
	K_Q      = C.SDLK_q
	K_F      = C.SDLK_f
	K_T      = C.SDLK_t
)

type Event interface{}

type EventQuit interface{}

type EventUnknown interface{}

type EventExposed interface{}

type EventMouseClick struct {
	X, Y int
	Down bool
}

type EventMouseWheel struct {
	X, Y int
}

type EventKey struct {
	Code int
	Down bool
}

type EventResize struct {
	W, H int
}

// Wait for next event
func WaitEvent() Event {
	var cev C.SDL_Event
	if 0 == C.SDL_WaitEvent(&cev) {
		return nil
	}
	return classifyEvent(&cev)
}

func WaitEventTimeout(timeout int) Event {
	var cev C.SDL_Event
	if 0 == C.SDL_WaitEventTimeout(&cev, C.int(timeout)) {
		return nil
	}
	return classifyEvent(&cev)
}

// Poll for pending envents. Return nil if there is no event available
func PollEvent() Event {
	var cev C.SDL_Event
	if 0 == C.SDL_PollEvent(&cev) {
		return nil
	}
	return classifyEvent(&cev)
}

func classifyEvent(cev *C.SDL_Event) Event {
	switch C.getEventType(cev) {

	case C.SDL_QUIT:
		return new(EventQuit)

	case C.SDL_KEYDOWN:
		// Ignore repeat key events
		if C.isKeyRepeat(cev) != 0 {
			break
		}
		kde := new(EventKey)
		kde.Code = int(C.getKeyCode(cev))
		kde.Down = true
		return kde

	case C.SDL_KEYUP:
		kde := new(EventKey)
		kde.Code = int(C.getKeyCode(cev))
		kde.Down = false
		return kde

	case C.SDL_WINDOWEVENT:
		wet := C.getWindowEvent(cev)
		if wet == C.SDL_WINDOWEVENT_RESIZED {
			var w, h C.int
			C.getWinData(cev, &w, &h)
			wr := new(EventResize)
			wr.W = int(w)
			wr.H = int(h)
			return wr

		} else if wet == C.SDL_WINDOWEVENT_EXPOSED {
			return new(EventExposed)
		}

	case C.SDL_MOUSEBUTTONDOWN:
		var x, y, dw C.int
		var down bool
		C.getMouseCoords(cev, &x, &y, &dw)
		if dw == 1 {
			down = true
		}
		return &EventMouseClick{X: int(x), Y: int(y), Down: down}

	case C.SDL_MOUSEWHEEL:
		var x, y C.int
		C.getMouseWheel(cev, &x, &y)
		return &EventMouseWheel{int(x), int(y)}
	}
	return new(EventUnknown)
}

// Process events. Returns true if EventQuit has appeared
func SlurpEvents() (quit bool) {
	quit = false
	for {
		ev := PollEvent()
		if ev == nil {
			return
		}
		switch ev.(type) {
		case *EventQuit:
			quit = true
		}
	}
}

// Returns true if key is pressed, false otherwise
func IsKeyPressed(kcode int) bool {
	return C.isKeyPressed(C.int(kcode)) == 1
}
