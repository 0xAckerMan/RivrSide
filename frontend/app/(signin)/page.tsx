"use client";
import SignInForm from "@/components/Signin";
import { ModeToggle } from "@/components/darktoggle";

const SignInPage = () => {
  return (
    <>
    <div className="flex flex-row-reverse pr-3">
      <ModeToggle/>
</div>
      <SignInForm />
    </>
  );
}

export default SignInPage;