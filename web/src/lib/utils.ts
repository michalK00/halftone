import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
interface Mime {
  mime: string,
  pattern: number[],
  mask: number[],
}

const mimes: Mime[] = [
  {
    mime: 'image/jpeg',
    pattern: [0xFF, 0xD8, 0xFF],
    mask: [0xFF, 0xFF, 0xFF],
  },
  {
    mime: 'image/png',
    pattern: [0x89, 0x50, 0x4E, 0x47],
    mask: [0xFF, 0xFF, 0xFF, 0xFF],
  }
  // you can expand this list @see https://mimesniff.spec.whatwg.org/#matching-an-image-type-pattern
];
export function getMimeType(file: File): Promise<string> {
  return new Promise((resolve) => {

    function check(bytes: Uint8Array, mime: Mime): boolean {
      for (let i = 0, l = mime.mask.length; i < l; ++i) {
        if ((bytes[i] & mime.mask[i]) - mime.pattern[i] !== 0) {
          return false;
        }
      }
      return true;
    }

    // Only read the first 4 bytes of the file
    const blob = file.slice(0, 4);
    const reader = new FileReader();

    reader.onloadend = function(e) {
      if (e.target?.readyState === FileReader.DONE && e.target.result) {
        const bytes = new Uint8Array(e.target.result as ArrayBuffer);

        for (const mime of mimes) {
          if (check(bytes, mime)) {
            resolve(mime.mime);
            return;
          }
        }
        resolve(file.type);
      }
    };

    reader.onerror = function() {
      resolve(file.type);
    };

    reader.readAsArrayBuffer(blob);
  });
}