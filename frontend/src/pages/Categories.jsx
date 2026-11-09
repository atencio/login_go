import CatalogPage from './CatalogPage';

const CONFIG = {
  path: '/categories', perm: 'categories', singular: 'categoría', example: 'Ej. Monitor',
  title: 'Categorías de equipos', subtitle: 'Tipos de equipo que se eligen al registrar un equipo.',
  activeLabel: 'Categoría activa (aparece al registrar equipos)',
  texts: {
    created: 'Categoría creada', updated: 'Categoría actualizada', deleted: 'Categoría eliminada',
    saveError: 'No se pudo guardar la categoría',
    createdText: 'Ya puedes elegirla al registrar equipos.', updatedText: 'Los cambios se guardaron.',
  },
};

export default function Categories() {
  return <CatalogPage cfg={CONFIG} />;
}
