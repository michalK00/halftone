import React, { createContext, useContext, useEffect, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api.ts"

const authApi = {
    signIn: async (email: string, password: string) => {
        console.log(import.meta.env.VITE_API_BACKEND_URL);
        const response = await api.post("/auth/signin", { email, password });
        return response.data;
    },

    signUp: async (userData: {  email: string; password: string }) => {
        const response = await api.post("/auth/signup", userData);
        return response.data;
    },

    verifyAccount: async (email: string, code: string) => {
        const response = await api.post("/auth/verify", { email, code });
        return response.data;
    },

    resendCode: async (username: string) => {
        const response = await api.post("/auth/resend-verification", { username });
        return response.data;
    },

    forgotPassword: async (username: string) => {
        const response = await api.post("/auth/forgot-password", { username });
        return response.data;
    },

    resetPassword: async (username: string, code: string, newPassword: string) => {
        const response = await api.post("/auth/reset-password", {
            username,
            code,
            newPassword
        });
        return response.data;
    },

    refreshToken: async () => {
        const refreshToken = localStorage.getItem("refreshToken");
        if (!refreshToken) throw new Error("No refresh token available");

        const response = await api.post("/auth/refresh-token", { refresh_token: refreshToken });
        return response.data;
    },
};

interface AuthContextType {
    isAuthenticated: boolean;
    isLoading: boolean;
    signIn: (email: string, password: string) => Promise<any>;
    signUp: (email: string, password: string) => Promise<any>;
    verifyAccount: (email: string, code: string) => Promise<any>;
    resendVerificationCode: (username: string) => Promise<any>;
    forgotPassword: (username: string) => Promise<any>;
    resetPassword: (username: string, code: string, newPassword: string) => Promise<any>;
    signOut: () => void;
    error: Error | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [isInitializing, setIsInitializing] = useState<boolean>(true);
    const [error, setError] = useState<Error | null>(null);
    const queryClient = useQueryClient();

    useEffect(() => {
        const checkAuthState = () => {
            const token = localStorage.getItem("token");

            if (token) {
                setIsAuthenticated(true);
            }
            setIsInitializing(false);
        };

        checkAuthState();
    }, []);

    const signInMutation = useMutation({
        mutationFn: ({ email, password }: { email: string; password: string }) =>
            authApi.signIn(email, password),
        onSuccess: (data) => {
            localStorage.setItem("token", data.id_token);
            localStorage.setItem("refreshToken", data.refresh_token);

            setIsAuthenticated(true);
            setError(null);
            queryClient.invalidateQueries({ queryKey: ["user"] });
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const signUpMutation = useMutation({
        mutationFn: ({ email, password }: { email: string; password: string }) =>
            authApi.signUp({ email, password }),
        onSuccess: () => {
            setError(null);
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const verifyAccountMutation = useMutation({
        mutationFn: ({ email, code }: { email: string; code: string }) =>
            authApi.verifyAccount(email, code),
        onSuccess: () => {
            setError(null);
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const resendCodeMutation = useMutation({
        mutationFn: (username: string) => authApi.resendCode(username),
        onSuccess: () => {
            setError(null);
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const forgotPasswordMutation = useMutation({
        mutationFn: (username: string) => authApi.forgotPassword(username),
        onSuccess: () => {
            setError(null);
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const resetPasswordMutation = useMutation({
        mutationFn: ({ username, code, newPassword }: { username: string; code: string; newPassword: string }) =>
            authApi.resetPassword(username, code, newPassword),
        onSuccess: () => {
            setError(null);
        },
        onError: (err: Error) => {
            setError(err);
        }
    });

    const signOut = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("refreshToken");
        setIsAuthenticated(false);
        queryClient.clear();
    };

    const value: AuthContextType = {
        isAuthenticated,
        isLoading: isInitializing || signInMutation.isPending,
        signIn: (email, password) => signInMutation.mutateAsync({ email, password }),
        signUp: (email, password) => signUpMutation.mutateAsync({ email, password }),
        verifyAccount: (email, code) => verifyAccountMutation.mutateAsync({ email, code }),
        resendVerificationCode: (username) => resendCodeMutation.mutateAsync(username),
        forgotPassword: (username) => forgotPasswordMutation.mutateAsync(username),
        resetPassword: (username, code, newPassword) => resetPasswordMutation.mutateAsync({ username, code, newPassword }),
        signOut,
        error
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};

// Export API client for use in other parts of the app
export { api };