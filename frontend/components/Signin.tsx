"use client";
import { FormEventHandler, useState, useEffect } from "react";
import { toast } from "sonner";

const SignInForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  // Function to handle form submission
  const handleSubmit: FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();

    try {
      setLoading(true);
      await signIn(email, password);
      toast.success("Successfully signed in!");
    } catch {
      toast.error("Wrong email or password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center h-[82vh]">
      <div className="md:m-auto md:w-96 space-y-9 w-auto m-auto">
        <div className="flex flex-col items-center">
          <h1 className="font-medium text-3xl">Welcome Back</h1>
          <p>Enter your details</p>
        </div>

        {/* <div
          className="mx-2 flex flex-col items-center"
          style={{ minHeight: "24px" }}
        >
          {error && <p className="text-red-500">{error}</p>}
        </div> */}
        <form onSubmit={handleSubmit} className="flex flex-col">
          <label className="font-medium mx-2">Email Address</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter your email"
            className="border-2 p-2 m-2 rounded-md active:outline-none focus:outline-none border-blue-300 focus:border-blue-500 hover:border-blue-500"
          />

          <label className="font-medium mx-2">Password</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password"
            className="border-2 p-2 m-2 rounded-md active:outline-none focus:outline-none border-blue-300 focus:border-blue-500 hover:border-blue-500"
          />
          <button
            disabled={loading}
            type="submit"
            className="border-2 border-[#076B84] rounded-md bg-[#076B84] m-2 p-2 text-white hover:bg-[#386792]"
          >
            LogIn
          </button>
        </form>
      </div>
    </div>
  );
};

export default SignInForm;
