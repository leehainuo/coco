/**
 * useIconProps - Unified icon component props
 * Provides standard size and color props for all icon components
 */
export function useIconProps() {
  return {
    size: { type: Number, default: 24 },
    color: { type: String, default: 'currentColor' },
  }
}
