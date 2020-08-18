import Czzle from './czzle';

Object.defineProperty(window, 'czzle', {
    get() {
      return Czzle;
    },
});