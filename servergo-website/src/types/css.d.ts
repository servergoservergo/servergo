declare module '*.css' {
  const content: { [className: string]: string };
  export default content;
}

declare module 'prismjs/themes/*';
declare module 'prismjs/plugins/*/prism-*';
declare module 'prismjs/plugins/*/*.css'; 