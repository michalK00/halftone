import { useState, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useToast } from "@/hooks/use-toast.ts";

const VerifyAccountPage = () => {
    const location = useLocation();
    const [username, setUsername] = useState("");
    const [code, setCode] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [isResending, setIsResending] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const { verifyAccount, resendVerificationCode } = useAuth();
    const navigate = useNavigate();
    const { toast } = useToast();

    useEffect(() => {
        // Get username from state if available
        if (location.state && location.state.username) {
            setUsername(location.state.username);
        }
    }, [location]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            await verifyAccount(username, code);

            toast({
                title: "Account verified",
                description: "Your account has been successfully verified. You can now sign in.",
                variant: "default",
            });

            navigate("/signin");
        } catch (err: any) {
            console.error("Verification error:", err);
            setError(err.message || "Failed to verify account. Please check the code and try again.");
        } finally {
            setIsLoading(false);
        }
    };

    const handleResendCode = async () => {
        if (!username) {
            setError("Username is required to resend verification code");
            return;
        }

        setIsResending(true);
        setError(null);

        try {
            await resendVerificationCode(username);

            toast({
                title: "Code resent",
                description: "A new verification code has been sent to your email",
                variant: "default",
            });
        } catch (err: any) {
            console.error("Resend code error:", err);
            setError(err.message || "Failed to resend verification code. Please try again.");
        } finally {
            setIsResending(false);
        }
    };

    return (
        <div className="flex justify-center items-center min-h-screen bg-background p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold">Verify your account</CardTitle>
                    <CardDescription>Enter the verification code sent to your email</CardDescription>
                </CardHeader>
                <CardContent>
                    {error && (
                        <Alert variant="destructive" className="mb-4">
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="username">Username</Label>
                            <Input
                                id="username"
                                type="text"
                                placeholder="Enter your username"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                required
                                disabled={isLoading || isResending}
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="code">Verification Code</Label>
                            <Input
                                id="code"
                                type="text"
                                placeholder="Enter verification code"
                                value={code}
                                onChange={(e) => setCode(e.target.value)}
                                required
                                disabled={isLoading || isResending}
                            />
                        </div>
                        <Button type="submit" className="w-full" disabled={isLoading || isResending}>
                            {isLoading ? (
                                <>
                                    <span className="animate-spin mr-2">⟳</span> Verifying...
                                </>
                            ) : (
                                "Verify Account"
                            )}
                        </Button>
                    </form>
                </CardContent>
                <CardFooter className="flex flex-col space-y-4">
                    <Button
                        variant="outline"
                        className="w-full"
                        onClick={handleResendCode}
                        disabled={isLoading || isResending}
                    >
                        {isResending ? (
                            <>
                                <span className="animate-spin mr-2">⟳</span> Resending...
                            </>
                        ) : (
                            "Resend Code"
                        )}
                    </Button>
                    <div className="text-sm text-center">
                        <Link to="/signin" className="text-primary hover:underline">
                            Back to Sign In
                        </Link>
                    </div>
                </CardFooter>
            </Card>
        </div>
    );
};

export default VerifyAccountPage;