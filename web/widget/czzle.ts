import { Config } from './config';
import { VDom, Component, VElement, ChangeDetection, SafeHTML, VStyleMap } from './vdom';

const dom = new VDom('czzle');

function configure(cfg: Config) {
    document.querySelectorAll('input[type="hidden"][czzle]').forEach((el) => {
        let name: string;
        if (el.hasAttribute('name')) {
            name = el.getAttribute('name') as string;
        } else {
            name = 'czzle';
        }
        dom.bootstrap(el as any, new CheckBox(cfg, name));
    })
    dom.bootstrap(document.head, new GlobalStyle(), false);
    window.addEventListener('load', () => {
        console.warn("loaded", cfg)
    });
}

export default {
    configure,
}

let itr = 0;
export class CheckBox implements Component {
    private id = itr++;
    private checked = false;
    private value = '';
    private onChange: ChangeDetection = () => {};
    constructor(
        private cfg: Config,
        private name: string,
    ) { }
    element(): VElement {
        return dom.create('div', {
            style: {
                'border': '1px solid rgb(184, 184, 184)',
                'max-width': '200px',
                'padding': '16px',
                'border-radius': '3px',
                'background-color': '#fafafa',
                'display': 'flex',
                'flex-direction': 'row',
                'align-items': 'center',
                'justify-content': 'start',
                'margin': '10px 0',
                'cursor': this.checked ? 'default' : 'pointer',
                '-webkit-touch-callout': 'none', 
                '-webkit-user-select': 'none', 
                '-khtml-user-select': 'none', 
                '-moz-user-select': 'none', 
                '-ms-user-select': 'none', 
                'user-select': 'none',
            },
            children: [
                dom.create('input', {
                    attrs: {
                        type: 'hidden',
                        name: this.name,
                        value: this.value,
                    }
                }),
                dom.create('div', {
                    style: {
                        'margin-right': '5px',
                        'height': '24px',
                        'width': '24px',
                    },
                    children: [
                        this.checked
                        ? new SafeHTML(
                            `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">`+
                                `<path d="M0 0h24v24H0z" fill="none"/>`+
                                `<path fill="#7b8794" d="M19 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.11 0 2-.9 2-2V5c0-1.1-.89-2-2-2zm-9 14l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>`+
                            `</svg>`
                        )
                        : new SafeHTML(
                            `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">`+
                                `<path fill="#9e9e9e" d="M19 5v14H5V5h14m0-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2z"/>`+
                                `<path d="M0 0h24v24H0z" fill="none"/> `+
                            `</svg>`
                        ),
                    ],
                }),
                dom.create('span', {
                    style: {
                        'height': '24px',
                        'line-height': '24px',
                        'font-size': '18px',
                    },
                    children: [
                        `I'm not a robot`,
                    ],
                }),
            ],
            on: {
                click: () => {
                    if (this.value) {
                        return;
                    }
                    let undo = [];
                    fetch(`${this.cfg.apiUrl}/v1/begin`, {
                        method: 'POST',
                        body: JSON.stringify({}),
                        headers: {
                            'Content-Type': 'application/json; charset=UTF-8',
                        }
                    }).then( (res) => {
                        if (res.ok) {
                            return res.json();
                        }
                        return Promise.reject(res);
                    }).then( (res) => {
                        const puzzle  = dom.bootstrap(
                            document.body,
                            new PuzzleContainer(res.puzzle, (token) => {
                                if (token) {
                                    this.checked = true;
                                    this.value = token;
                                    this.onChange();
                                }
                                puzzle.remove();
                            }),
                            false,
                        )
                    });
                    
                }
            }
        });
    }
    registerChange(changecb: ChangeDetection) {
        this.onChange = changecb;
    }
}

export class GlobalStyle implements Component {
    element(): VElement {
        return dom.create('style', {
            children: [globalStyle],
        })
    }
}

export enum Level {
    Unknown = 'unknown',
    None = 'none',
    Easy = 'easy',
    Medium = 'medium',
    Hard = 'hard'
}

export interface ClientInfo {
    id: string;
    ip: string;
    user_agent: string;
    time: number;
}
export interface Puzzle {
    level: Level;
    token: string;
    tile_map: TileMap;
    expires_at: number;
    issued_at: number;
}

export interface Pos  {
    x: number;
    y: number;
}

export enum TileType {
    Front = 'front',
    Back = 'back',
}


export interface TileMap {
    size: number;
    tiles: Tile[];
}

export interface Tile {
    type: TileType;
    pos: Pos;
    data: string;
}
const globalStyle = `
    @keyframes czzle-puzzle-host{
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
    @keyframes czzle-flip-1 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(90deg);
        }
    }
    @keyframes czzle-flip-2 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(90deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
    }
    @keyframes czzle-flip-back-1 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(-90deg);
        }
    }
    @keyframes czzle-flip-back-2 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(-90deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
    }
`

export type SolveCb = (token?: string) => void;

export class PuzzleContainer implements Component {
    private onChange?: ChangeDetection;
    private grid: TileGrid;
    constructor(
        private puzzle: Puzzle,
        private slovecb: SolveCb,
    ) {
        this.grid = new TileGrid(this.puzzle.tile_map);
    }
    element(): VElement {
        return dom.create('div', {
            style: {
                'position': 'fixed',
                'display': 'flex',
                'flex-direction': 'row',
                'align-items': 'center',
                'justify-content': 'center',
                'top': 0,
                'left': 0,
                'width': '100vw',
                'height': '100vh',
                'background-color': 'rgba(0, 0, 0, 0.5)',
                'z-index': 99999999,
                'animation-name': 'czzle-puzzle-host',
                'animation-duration': '0.5s',
            },
            on: {
                click: () => {
                    // this.slovecb('trueasdas');
                },
            },
            children: [
                dom.create('div', {
                    style: {
                        'padding': '20px',
                        'background-color': '#fefefe',
                    },
                    children: [
                        this.grid.element(),
                    ],
                })
            ]
        })
    }
    registerChange(changecb: ChangeDetection) {
        this.onChange = changecb;
        this.grid.registerChange(changecb);
    }
}



export class TileGrid implements Component {
    private onChange?: ChangeDetection;
    private selected?: Pos | null;
    private animate: Pos[] = [];
    private animateBack: Pos[] = [];
    private animatePrev: Pos[] = [];
    private animatePrevBack: Pos[] = [];
    constructor(
        private m: TileMap,
    ) {} 
    element(): VElement {
        const tiles = this.m.tiles.filter(t => {
            if (this.selected) {
                if (t.type === TileType.Front) {
                    if (t.pos.x !== this.selected.x || t.pos.y !== this.selected.y) {
                       return true;
                    } 
                    return false;
                } else {
                    if (t.pos.x === this.selected.x && t.pos.y === this.selected.y) {
                        return true;
                    }
                    return false;
                }
            }
            return t.type === TileType.Front;
        });
        return dom.create('div', {
            style: {

                'height': '300px',
                'width': '300px',
                'position': 'relative',
            },
            children: tiles.map(t => {
                let extra: VStyleMap = {};
                if (this.animate.find(p => p.x === t.pos.x && p.y === t.pos.y)){
                    extra = {
                        'animation-name': 'czzle-flip-1',
                        'animation-duration': '0.25s',
                    }
                }
                if (this.animateBack.find(p => p.x === t.pos.x && p.y === t.pos.y)){
                    extra = {
                        'animation-name': 'czzle-flip-2',
                        'animation-duration': '0.25s',
                    }
                }
                if (this.animatePrev.find(p => p.x === t.pos.x && p.y === t.pos.y)){
                    extra = {
                        'animation-name': 'czzle-flip-back-1',
                        'animation-duration': '0.25s',
                    }
                }
                if (this.animatePrevBack.find(p => p.x === t.pos.x && p.y === t.pos.y)){
                    extra = {
                        'animation-name': 'czzle-flip-back-2',
                        'animation-duration': '0.25s',
                    }
                }
                if (this.selected) {
                    if (this.selected.x === t.pos.x && this.selected.y === t.pos.y) {
                        extra = {
                            'border': '5px solid rgb(23, 100, 163)',
                            'margin': '-5px 0 0 -5px',
                            ...extra,
                        }
                    }
                }
                return dom.create('img', {
                    attrs: {
                        'src': `data:image/png;base64,${t.data}`,
                    },
                    style: {
                        'position': 'absolute',
                        'left': `${t.pos.x * 100}px`,
                        'top': `${t.pos.y * 100}px`,
                        'border': 'none',
                        'width': '100px',
                        'height': '100px',
                        'outline': 'none',
                        'cursor': 'pointer',
                        ...extra,
                    },
                    on: {
                        click: () => {
                            if (this.animate.length || this.animateBack.length || this.animatePrev.length || this.animatePrevBack.length) {
                                return;
                            }
                            let selected: Pos | null;
                            if (this.selected) {
                                this.animatePrev = [this.selected];
                                if (this.selected.x === t.pos.x && this.selected.y === t.pos.y) {
                                    selected = null;
                                } else {
                                    selected = t.pos;
                                    this.animate = [t.pos];
                                }
                            }  else {
                                selected = t.pos;
                                this.animate = [t.pos];
                            }
                            setTimeout(() => {
                                this.selected = selected;
                                this.animateBack = this.animate;
                                this.animate = [];
                                this.animatePrevBack = this.animatePrev;
                                this.animatePrev = [];

                                setTimeout(() => {
                                    this.animateBack = [];
                                    this.animatePrevBack = [];
                                    if (this.onChange) {
                                        this.onChange();   
                                    }
                                }, 250)
                                if (this.onChange) {
                                    this.onChange();   
                                }
                            }, 250);
                            if (this.onChange) {
                                this.onChange();   
                            }
                        }
                    }
                });
            }),
        })
    }
    registerChange(changecb: ChangeDetection) {
        this.onChange = changecb;
    }
}