export interface VAttributeMap {
    [key: string]: string | number | boolean;
}



export interface VStyleMap {
    [key: string]: string | number;
}

export type VNode = VElement | string | SafeHTML;

export interface VElement {
    tag: string;
    on: VEventMap;
    attrs: VAttributeMap;
    style: VStyleMap;
    children: VNode[];
}

export interface VInputChangeEvent {
    value: boolean | number | string;
}

export interface VMouseClickEvent {
    x: number;
    y: number;
}

export interface VFormSubmitEvent {
    values: { [k: string]: string | boolean | number };
}

export interface VEventMap {
    click?: (e: VMouseClickEvent) => void;
    submit?: (e: VFormSubmitEvent) => void;
    change?: (e: VInputChangeEvent) => void;
}

export interface CreateOptions {
    on?: VEventMap;
    attrs?: VAttributeMap;
    style?: VStyleMap;
    children?: VNode[];
}

export type Patch = (node: HTMLElement) => any;

export class SafeHTML {
    constructor(public innerHTML: string) { }
}

export class VDom {
    private uid: string;
    constructor(uid: string) {
        this.uid = uid;
    }
    public create(
        tag: string,
        { on, attrs, style, children }: CreateOptions = {}
    ): VElement {
        if (attrs) {
            for (const k in attrs) {
                let v = attrs[k];
                if (k === 'id' || k === 'class' || k === 'for') {
                    v = this.prefixAttr(v);
                    attrs[k] = v;
                }
            }
        }
        return {
            tag: tag,
            style: style ? style : {},
            on: on ? on : {},
            attrs: attrs ? attrs : {},
            children: children ? children : []
        };
    }
    private prefixAttr(value: string | boolean | number): string {
        const list = value.toString().split(' ');
        for (let i = 0; i < list.length; i++) {
            list[i] = `${this.uid}-${list[i]}`;
        }
        return list.join(' ');
    }
    public renderElement(el: VElement): Node {
        const $el = document.createElement(el.tag);
        for (const k in el.attrs) {
            const v = el.attrs[k];
            $el.setAttribute(k, v.toString());
        }
        for (const child of el.children) {
            const $child = this.render(child);
            $el.appendChild($child);
        }
        let styles = [];
        for (const k in el.style) {
            let v = el.style[k];
            if (typeof v === 'number') {
                v = `${v}px`;
            }
            styles.push(`${k}: ${el.style[k]}`);
        }
        if (styles.length) {
            $el.setAttribute('style', styles.join('; '));
        }
        if (el.on.click) {
            const cb = el.on.click;
            $el.onclick = e => {
                cb({
                    x: e.x,
                    y: e.y
                });
            };
        }
        if (el.on.change) {
            const cb = el.on.change;
            const type = $el.getAttribute('type');
            if (type) {
                $el.onchange = e => {
                    if (type === 'checkbox') {
                        cb({
                            value: ($el as any)['checked']
                        });
                    } else {
                        cb({
                            value: ($el as any)['value']
                        });
                    }
                };
            }
        }
        return $el;
    }
    public renderText(text: string): Node {
        return document.createTextNode(text);
    }
    public renderHTML(node: SafeHTML): Node {
        const template = document.createElement('template');
        const html = node.innerHTML.trim();
        template.innerHTML = html;
        return template.content.firstChild as any;
    }
    public render(node: VNode): Node {
        if (typeof node === 'string') {
            return this.renderText(node);
        }
        if (node instanceof SafeHTML) {
            return this.renderHTML(node);
        }
        return this.renderElement(node);
    }
    private zip(xs: any, ys: any) {
        const zipped = [];
        for (let i = 0; i < Math.max(xs.length, ys.length); i++) {
            zipped.push([xs[i], ys[i]]);
        }
        return zipped;
    }

    private diffAttrs(oldAttrs: VAttributeMap, newAttrs: VAttributeMap): Patch {
        const patches: Patch[] = [];
        for (const k in newAttrs) {
            const v = newAttrs[k];
            patches.push($node => {
                $node.setAttribute(k, v.toString());
                return $node;
            });
        }
        for (const k in oldAttrs) {
            if (!(k in newAttrs)) {
                patches.push($node => {
                    $node.removeAttribute(k);
                    return $node;
                });
            }
        }
        return $node => {
            for (const patch of patches) {
                patch($node);
            }
        };
    }
    private diffStyle(oldStyle: VStyleMap, newStyle: VStyleMap): Patch {
        return $node => {
            let styles = [];
            for (const k in newStyle) {
                let v = newStyle[k];
                if (typeof v === 'number') {
                    v = `${v}px`;
                }
                styles.push(`${k}: ${newStyle[k]}`);
            }
            if (styles.length) {
                $node.setAttribute('style', styles.join('; '));
            }
        };
    }
    private diffEvents(oldMap: VEventMap, newMap: VEventMap): Patch {
        const patches: Patch[] = [];
        if (newMap.click) {
            const cb = newMap.click;
            patches.push($node => {
                $node.onclick = e => {
                    cb({
                        x: e.x,
                        y: e.y
                    });
                };
                return $node;
            });
        } else if (oldMap.click) {
            patches.push($node => {
                $node.onclick = null;
            });
        }
        if (newMap.change) {
            const cb = newMap.change;
            patches.push($node => {
                const type = $node.getAttribute('type');
                if (type) {
                    $node.onchange = e => {
                        if (type === 'checkbox') {
                            cb({
                                value: ($node as any)['checked']
                            });
                        } else {
                            cb({
                                value: ($node as any)['value']
                            });
                        }
                    };
                }
            });
        } else if (oldMap.change) {
            patches.push($node => {
                $node.onchange = null;
            });
        }
        return $node => {
            for (const patch of patches) {
                patch($node);
            }
        };
    }
    private diffChildren(oldVChildren: VNode[], newVChildren: VNode[]): Patch {
        const childPatches: Patch[] = [];
        oldVChildren.forEach((oldVChild, i) => {
            childPatches.push(this.diff(oldVChild, newVChildren[i]));
        });
        const additionalPatches: Patch[] = [];
        for (const additionalVChild of newVChildren.slice(
            oldVChildren.length
        )) {
            additionalPatches.push($node => {
                $node.appendChild(this.render(additionalVChild));
                return $node;
            });
        }
        return $parent => {
            for (const [patch, child] of this.zip(
                childPatches,
                $parent.childNodes
            )) {
                patch(child);
            }
            for (const patch of additionalPatches) {
                patch($parent);
            }
            return $parent;
        };
    }
    private diff(vOldNode: VNode, vNewNode: VNode): Patch {
        if (vNewNode === undefined) {
            return $node => {
                $node.remove();
                return undefined;
            };
        }
        if (typeof vOldNode === 'string' || typeof vNewNode === 'string') {
            if (vOldNode !== vNewNode) {
                return $node => {
                    const $newNode = this.render(vNewNode);
                    ($node as any).replaceWith($newNode);
                    return $newNode;
                };
            } else {
                return $node => undefined;
            }
        }
        if (vOldNode instanceof SafeHTML && vNewNode instanceof SafeHTML) {
            if (vOldNode.innerHTML !== vNewNode.innerHTML) {
                return $node => {
                    const $newNode = this.render(vNewNode);
                    ($node as any).replaceWith($newNode);
                    return $newNode;
                };
            } else {
                return $node => undefined;
            }
        } else if (
            vOldNode instanceof SafeHTML ||
            vNewNode instanceof SafeHTML
        ) {
            return $node => {
                const $newNode = this.render(vNewNode);
                ($node as any).replaceWith($newNode);
                return $newNode;
            };
        }

        if (vOldNode.tag !== vNewNode.tag) {
            return $node => {
                const $newNode = this.render(vNewNode);
                ($node as any).replaceWith($newNode);
                return $newNode;
            };
        }

        const patchAttrs = this.diffAttrs(vOldNode.attrs, vNewNode.attrs);
        const patchStyle = this.diffStyle(vOldNode.style, vNewNode.style);
        const patchEvents = this.diffEvents(vOldNode.on, vNewNode.on);
        const patchChildren = this.diffChildren(
            vOldNode.children,
            vNewNode.children
        );

        return $node => {
            patchAttrs($node);
            patchStyle($node);
            patchChildren($node);
            patchEvents($node);
            return $node;
        };
    }
    private mount(node: any, target: any): any {
        (target as any).replaceWith(node);
        return node;
    }
    public bootstrap(el: HTMLElement, cmp: Component, replace = true): HTMLElement {
        let app = cmp.element();
        let $app = document.createElement('div');
        if (replace) {
            el.replaceWith($app);
        } else {
            el.append($app);
        }
        $app = this.mount(this.render(app), $app);
        if (cmp.registerChange) {
            cmp.registerChange(() => {
                const newApp = cmp.element();
                const patch = this.diff(app, newApp);
                $app = patch($app);
                app = newApp;
            });
        }
        return $app;
    }
}


export type ChangeDetection = () => void;

export interface Component {
    element(): VElement;
    registerChange?(cb: ChangeDetection): void;
}
