import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';
import { jwtVerify } from 'jose';
import ProfileClient from './ProfileClient';

export default async function ProfilePage() {
    const cookieStore = cookies();
    const token = cookieStore.get('session_token')?.value;
    
    if (!token) {
        console.log("meow");
        redirect('/login');
    }

    try {
        await jwtVerify(token, new TextEncoder().encode(process.env.JWT_PRIVATE_KEY));
        console.log("awaited")
    } catch (e) {
        console.error("JWT Verify Failed", e);
        redirect('/login');
    }

    return <ProfileClient/>;
}
