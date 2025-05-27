declare module '@uidotdev/usehooks' {
  export function useLocalStorage<T>(key: string, initialValue: T): [T, (value: T) => void];
  export function useMediaQuery(query: string): boolean;
  export function useOnClickOutside<T extends HTMLElement>(
    ref: React.RefObject<T>,
    handler: (event: MouseEvent | TouchEvent) => void
  ): void;
  export function useDebounce<T>(value: T, delay: number): T;
  export function useWindowSize(): {
    width: number;
    height: number;
  };
} 