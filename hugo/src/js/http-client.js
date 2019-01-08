import axios from 'axios';
import { cacheAdapterEnhancer, throttleAdapterEnhancer } from 'axios-extensions';

export default axios.create({
  adapter: throttleAdapterEnhancer(cacheAdapterEnhancer(axios.defaults.adapter)),
});
