import * as React from "react";
import { Outlet } from "react-router-dom";

import { useLocation, useNavigate, useSearchParams } from "react-router-dom";
import { twMerge } from "tailwind-merge";

import { toast } from "react-hot-toast";
import { isDev } from "@/utils";
import * as Auth from "@/apis/authv1/authv1";
import { getClientAuth } from "@/utils/client";
import { useMutation } from "@tanstack/react-query";
import { Divider } from "@mantine/core";

interface authResponse {
  loginURL: string;
}

interface authReqCommon {
  query?: string;
  userAgent: string;
}

interface State {
  domain: string;
  identityProviders?: StateProvider[];
  isPasskeyLoginEnabled?: boolean;
}

interface StateProvider {
  uid: string;
  displayName: string;
  picURL?: string;
}

function getState() {
  if (!isDev()) {
    return (window as any).__OCTELIUM_STATE__ as State;
  }

  return {
    domain: "example.com",
    isPasskeyLoginEnabled: true,
    identityProviders: [
      {
        uid: "github",
        displayName: "GitHub",
      },
      {
        uid: "gitlab-1",
        displayName: "Gitlab",
      },
      {
        uid: "gitlab-2",
        displayName: "Gitlab",
      },
    ],
  } as State;
}

const Passkey = () => {
  const c = getClientAuth();

  const mutation = useMutation({
    mutationFn: async () => {
      const { response } = await c.authenticateWithPasskeyBegin(
        Auth.AuthenticateWithPasskeyBeginRequest.create({})
      );

      try {
        console.log("Got req", response.request);
        const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(
          JSON.parse(response.request)
        );
        const credential = (await navigator.credentials.get({
          publicKey,
          mediation: "conditional",
        })) as PublicKeyCredential;

        console.log("Got credential", credential.toJSON());

        return await c.authenticateWithPasskey(
          Auth.AuthenticateWithPasskeyRequest.create({
            response: JSON.stringify(credential.toJSON()),
          })
        );
      } catch (err) {
        console.log("fido get err", err);
        throw err;
      }
    },
    onSuccess: (r) => {
      window.location.href = "/callback/success";
    },
    onError: (resp) => {},
  });

  return (
    <div className="w-full">
      <button
        disabled={mutation.isPending}
        className={twMerge(
          "w-full px-2 py-4 md:py-6 font-bold transition-all duration-500 mb-4",
          "shadow-2xl rounded-lg cursor-pointer font-bold",
          "bg-[#242323] hover:bg-black text-white text-lg",
          mutation.isPending ? "!bg-[#777] shadow-none" : undefined
        )}
        onClick={() => {
          mutation.mutate();
        }}
      >
        <span className="font-bold text-lg">Login with Passkey</span>
      </button>
    </div>
  );
};

const Page = () => {
  const state = getState();

  let [loginActive, setLoginActive] = React.useState<boolean>(false);
  let [reqCommon, setReqCommon] = React.useState<authReqCommon | null>(null);

  const [searchParams, setSearchParams] = useSearchParams();

  React.useEffect(() => {
    setReqCommon({
      query: searchParams.toString() ?? "",
      userAgent: window.navigator.userAgent,
    });

    if (searchParams.has("error")) {
      toast.error(searchParams.get("error"));
    }
    searchParams.forEach((val, key, parent) => {
      searchParams.delete(key);
    });
    setSearchParams(searchParams);
  }, []);

  return (
    <div>
      <div className="flex items-center justify-center mt-4 mb-3">
        <div className="w-40 h-40 md:w-60 md:h-60 rounded-[50%] bg-black shadow-2xl flex items-center justify-center transition-all duration-300 hover:bg-zinc-900 hover:shadow-xl">
          <svg
            className="w-20 h-20 md:w-40 md:h-40"
            width="256"
            height="256"
            viewBox="0 0 256 256"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M138.806 201.044C123.741 197.565 108.138 197.087 92.8878 199.637C77.6378 202.188 63.0395 207.717 49.9263 215.91L74.9591 255.978C82.8103 251.073 91.5508 247.763 100.681 246.236C109.812 244.709 119.154 244.995 128.174 247.078L138.806 201.044Z"
              fill="white"
            />
            <path
              d="M187.291 172.009C174.178 180.201 162.807 190.896 153.827 203.483C144.847 216.07 138.435 230.303 134.955 245.368L180.989 256C183.072 246.98 186.912 238.459 192.288 230.922C197.665 223.386 204.473 216.983 212.324 212.078L187.291 172.009Z"
              fill="white"
            />
            <path
              d="M201.044 117.194C197.565 132.259 197.087 147.862 199.637 163.112C202.188 178.362 207.717 192.961 215.91 206.074L255.978 181.041C251.073 173.19 247.763 164.449 246.236 155.319C244.709 146.188 244.995 136.846 247.078 127.826L201.044 117.194Z"
              fill="white"
            />
            <path
              d="M172.009 68.7086C180.201 81.8217 190.896 93.1929 203.483 102.173C216.07 111.153 230.303 117.565 245.368 121.045L256 75.0112C246.98 72.9279 238.459 69.0884 230.922 63.7119C223.386 58.3354 216.983 51.5271 212.078 43.6759L172.009 68.7086Z"
              fill="white"
            />
            <path
              d="M117.194 54.9556C132.259 58.4351 147.862 58.9132 163.112 56.3626C178.362 53.812 192.961 48.2827 206.074 40.0903L181.041 0.0216198C173.19 4.92665 164.449 8.23722 155.319 9.76434C146.188 11.2915 136.846 11.0052 127.826 8.92191L117.194 54.9556Z"
              fill="white"
            />
            <path
              d="M68.7086 83.9909C81.8217 75.7986 93.1928 65.1036 102.173 52.5166C111.153 39.9297 117.565 25.6973 121.045 10.632L75.0112 0C72.9279 9.02004 69.0884 17.5414 63.7119 25.0776C58.3353 32.6138 51.5271 39.0172 43.6759 43.9222L68.7086 83.9909Z"
              fill="white"
            />
            <path
              d="M54.9556 138.806C58.4351 123.741 58.9132 108.138 56.3626 92.8879C53.812 77.6378 48.2827 63.0395 40.0903 49.9264L0.0216079 74.9591C4.92663 82.8103 8.23721 91.5508 9.76433 100.681C11.2914 109.812 11.0052 119.154 8.9219 128.174L54.9556 138.806Z"
              fill="white"
            />
            <path
              d="M83.9909 187.291C75.7986 174.178 65.1036 162.807 52.5166 153.827C39.9297 144.847 25.6973 138.435 10.632 134.955L0 180.989C9.02004 183.072 17.5414 186.912 25.0776 192.288C32.6138 197.665 39.0172 204.473 43.9222 212.324L83.9909 187.291Z"
              fill="white"
            />
          </svg>
        </div>
      </div>

      {(!state.identityProviders || state.identityProviders.length < 1) && (
        <div className="container mx-auto mt-2 p-2 md:p-8 w-full max-w-lg">
          {!state.isPasskeyLoginEnabled && (
            <h2 className="font-bold text-2xl text-slate-700 flex items-center justify-center mb-4 text-center">
              No Available Identity Providers
            </h2>
          )}

          {state.isPasskeyLoginEnabled && (
            <div>
              <Passkey />
            </div>
          )}
        </div>
      )}
      {state.identityProviders && state.identityProviders.length > 0 && (
        <div className="container mx-auto mt-2 p-2 md:p-4 w-full max-w-lg">
          <div
            className="font-bold text-xl mb-4 text-zinc-700 text-center"
            style={{
              textShadow: "0 2px 8px rgba(0, 0, 0, 0.2)",
            }}
          >
            <span>Login to</span>
            <span> </span>
            <span className="text-black">Octelium</span>
            <span> </span>
            <span>with an Identity Provider</span>
          </div>

          <div className="flex flex-col items-center justify-center">
            {state.identityProviders.map((c) => {
              return (
                <button
                  className={twMerge(
                    "w-full px-2 py-4 md:py-6 font-bold transition-all duration-500 mb-4",
                    "shadow-2xl rounded-lg cursor-pointer font-bold",
                    "bg-[#242323] hover:bg-black text-white text-lg",
                    loginActive ? "!bg-[#777] shadow-none" : undefined
                  )}
                  disabled={loginActive}
                  key={c.uid}
                  onClick={() => {
                    setLoginActive(true);
                    fetch("/begin", {
                      method: "POST",
                      headers: {
                        "Content-Type": "application/json",
                        Accept: "application/json",
                      },
                      body: JSON.stringify({
                        uid: c.uid,
                        ...reqCommon,
                      }),
                    })
                      .then((res) => res.json())
                      .then((data: authResponse) => {
                        window.location.href = data.loginURL;
                      });
                  }}
                >
                  <div className="w-full flex flex-row items-center justify-center">
                    <span className="flex-1 flex items-center justify-center font-semibold">
                      {c.displayName}
                    </span>
                  </div>
                </button>
              );
            })}
          </div>
          {state.isPasskeyLoginEnabled && (
            <div>
              <Divider my="lg" label="OR" labelPosition="center" />
              <Passkey />
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Page;
