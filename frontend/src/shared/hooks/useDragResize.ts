/**
 * Drag resize composable
 */

import type { Ref } from 'vue'

export interface DragResizeOptions {
  /** Drag axis */
  axis: 'x' | 'y'
  /** Minimum value */
  min: number
  /** Maximum value */
  max: number
  /** Direction: 1 for right/down increase, -1 for right/down decrease */
  sign?: 1 | -1
}

/**
 * Start drag resize
 */
export function startDragResize(event: MouseEvent, sizeRef: Ref<number>, options: DragResizeOptions): void {
  const { axis, min, max, sign = 1 } = options
  const start = axis === 'x' ? event.clientX : event.clientY
  const val = sizeRef.value

  function onMove(ev: MouseEvent): void {
    const delta = (axis === 'x' ? ev.clientX : ev.clientY) - start
    sizeRef.value = Math.min(max, Math.max(min, val + sign * delta))
  }

  function onUp(): void {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
  }

  document.body.style.cursor = axis === 'x' ? 'col-resize' : 'row-resize'
  document.body.style.userSelect = 'none'
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
}
