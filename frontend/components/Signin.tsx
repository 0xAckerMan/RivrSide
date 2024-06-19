"use client";
import { FormEventHandler, useState } from "react";

const SignInForm = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
    
        try {
            setError('');
            setLoading(true);
            await signIn(email, password);
        } catch {
            setError('Wrong email or password');
        }
    
        setLoading(false);
    };
    
    return (
        <div className="w-96 space-y-3 flex flex-col items-center ju">
            <div className=" flex flex-col items-center">
                <h1 className=" font-bold text-3xl">Welcome Back</h1>
                <p>Enter your details</p>
            </div>
            <div className="mx-2" style={{ minHeight: '24px' }}>
                {error && <p className="text-red-500">{error}</p>}
            </div>
            <form onSubmit={handleSubmit} className="flex flex-col">
                <label className=" font-medium mx-2">Email Address</label>
                <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Enter your email"
                    className="border-2 border-gray-200 p-2 m-2 rounded-md active:outline-none focus:outline-none focus:border-blue-500"
                />

                <label className=" font-medium mx-2">Password</label>
                <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Password"
                    className="border-2 border-gray-200 p-2 m-2 rounded-md active:outline-none focus:outline-none focus:border-blue-500"
                />
                <button disabled={loading} type="submit" className="border-2 border-[#076B84] rounded-md bg-[#076B84] m-2 p-2 text-white hover:bg-[#386792]">
                    LogIn
                </button>
            </form>
        </div>
    );
}

export default SignInForm;
