const PATHS = {
  dashboard: 'M4 13h6V4H4v9zM14 20h6v-9h-6v9zM4 20h6v-4H4v4zM14 8h6V4h-6v4z',
  users: 'M16 19v-1a4 4 0 00-4-4H7a4 4 0 00-4 4v1M9.5 10a3.5 3.5 0 100-7 3.5 3.5 0 000 7zM21 19v-1a4 4 0 00-3-3.9M16 3.1a3.5 3.5 0 010 6.8',
  shield: 'M12 3l8 3v6c0 4.5-3.2 8-8 9-4.8-1-8-4.5-8-9V6l8-3zM9 12l2 2 4-4',
  chart: 'M4 20V10M10 20V4M16 20v-7M22 20H2',
  logout: 'M15 17l5-5-5-5M20 12H9M12 20H5a1 1 0 01-1-1V5a1 1 0 011-1h7',
  chevron: 'M15 6l-6 6 6 6',
  chevronDown: 'M6 9l6 6 6-6',
  menu: 'M4 7h16M4 12h16M4 17h16',
  search: 'M11 18a7 7 0 100-14 7 7 0 000 14zM20 20l-4-4',
  expand: 'M4 9V4h5M20 9V4h-5M4 15v5h5M20 15v5h-5',
  check: 'M5 13l4 4L19 7',
  x: 'M6 6l12 12M18 6L6 18',
  alert: 'M12 8v5M12 17h.01M10.3 3.9L2.4 18a2 2 0 001.7 3h15.8a2 2 0 001.7-3L13.7 3.9a2 2 0 00-3.4 0z',
  info: 'M12 11v6M12 7h.01M12 21a9 9 0 100-18 9 9 0 000 18z',
  plus: 'M12 5v14M5 12h14',
  edit: 'M4 20h4L19 9l-4-4L4 16v4zM14 6l4 4',
  trash: 'M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3',
  eye: 'M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12zM12 15a3 3 0 100-6 3 3 0 000 6z',
  eyeOff: 'M3 3l18 18M10.6 6.1A9.6 9.6 0 0112 6c6.4 0 10 6 10 6a17 17 0 01-3.2 3.9M6.6 7.6A17 17 0 002 12s3.6 6 10 6c1.5 0 2.8-.3 4-.8M9.9 9.9a3 3 0 004.2 4.2',
  refresh: 'M20 11a8 8 0 10-2.3 5.7M20 4v7h-7',
  user: 'M12 12a4 4 0 100-8 4 4 0 000 8zM4 20a8 8 0 0116 0',
};

export default function Icon({ name, size = 20 }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor"
         strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d={PATHS[name]} />
    </svg>
  );
}
