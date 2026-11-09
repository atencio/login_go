import CatalogPage from './CatalogPage';

const CONFIG = {
  path: '/brands', perm: 'brands', singular: 'marca', example: 'Ej. Lenovo',
  title: 'Marcas de equipos', subtitle: 'Fabricantes que se eligen al registrar un equipo.',
  activeLabel: 'Marca activa (aparece al registrar equipos)',
  texts: {
    created: 'Marca creada', updated: 'Marca actualizada', deleted: 'Marca eliminada',
    saveError: 'No se pudo guardar la marca',
    createdText: 'Ya puedes elegirla al registrar equipos.', updatedText: 'Los cambios se guardaron.',
  },
};

export default function Brands() {
  return <CatalogPage cfg={CONFIG} />;
}
