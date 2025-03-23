import { getUserIdFromLocalStorage, showError } from './utils';
import axios from 'axios';

const isDev =  process.env.NODE_ENV == 'development';
export let API = axios.create({
  baseURL: isDev? "": process.env.VITE_REACT_APP_SERVER_URL,
  withCredentials: true,
  headers: {
    'New-API-User': getUserIdFromLocalStorage(),
    'Cache-Control': 'no-store'
  }
});
export function updateAPI() {
  API = axios.create({
    baseURL: isDev? "": process.env.VITE_REACT_APP_SERVER_URL,
    withCredentials: true,
    headers: {
      'New-API-User': getUserIdFromLocalStorage(),
      'Cache-Control': 'no-store'
    }
  });
}

API.interceptors.response.use(
  (response) => response,
  (error) => {
    showError(error);
  },
);
