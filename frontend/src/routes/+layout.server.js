// import { redirect } from '@sveltejs/kit';

// export async function load({ cookies, url }) {
//     const token = cookies.get('token');
//     const isLoginPage = url.pathname.startsWith('/login');

//     if (!token && !isLoginPage) {
//         throw redirect(302, `/login?redirect=${url.pathname}`);
//     }

//     if (token && isLoginPage) {
//         throw redirect(302, '/');
//     }

//     return {};
// }
export const ssr = false;
