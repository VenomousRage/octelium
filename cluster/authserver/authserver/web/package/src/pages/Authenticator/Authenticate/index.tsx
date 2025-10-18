import OtpInput from "react-otp-input";
import * as React from "react";

import { isDev, queryClient } from "@/utils";

import * as Auth from "@/apis/authv1/authv1";
import { getClientAuth } from "@/utils/client";
import { useMutation, useQuery } from "@tanstack/react-query";
import TimeAgo from "@/components/TimeAgo";
import { getResourceRef } from "@/utils/pb";
import {
  ActionIcon,
  Button,
  Collapse,
  Input,
  Modal,
  PinInput,
  Tooltip,
} from "@mantine/core";
import { Timestamp } from "@/apis/google/protobuf/timestamp";
import { twMerge } from "tailwind-merge";
import { useDisclosure } from "@mantine/hooks";
import { DeleteOptions } from "@/apis/metav1/metav1";
import { MdEdit } from "react-icons/md";
import { MdEditOff } from "react-icons/md";
import { IoMdSend } from "react-icons/io";

const TOTP = (props: { authn: Auth.Authenticator }) => {
  const { authn } = props;

  const c = getClientAuth();

  const mutation = useMutation({
    mutationFn: async (otp: string) => {
      await c.authenticateAuthenticatorBegin(
        Auth.AuthenticateAuthenticatorBeginRequest.create({
          authenticatorRef: getResourceRef(authn),
        })
      );

      return await c.authenticateWithAuthenticator(
        Auth.AuthenticateWithAuthenticatorRequest.create({
          authenticatorRef: getResourceRef(authn),
          challengeResponse: {
            type: {
              oneofKind: "totp",
              totp: {
                response: otp,
              },
            },
          },
        })
      );
    },
    onSuccess: (r) => {
      window.location.href = "/callback/success";
    },
    onError: (resp) => {},
  });

  return (
    <div className="w-full flex flex-col items-center justify-center">
      <h2 className="font-bold text-xl text-slate-700 flex items-center justify-center my-4 text-center">
        Enter the OTP
      </h2>

      <PinInput
        type={"number"}
        inputType="tel"
        inputMode="numeric"
        length={6}
        autoFocus
        size="lg"
        onComplete={(val) => {
          mutation.mutate(val);
        }}
      />
    </div>
  );
};

const Fido = (props: { authn: Auth.Authenticator }) => {
  const { authn } = props;

  const c = getClientAuth();

  const mutation = useMutation({
    mutationFn: async () => {
      const { response } = await c.authenticateAuthenticatorBegin(
        Auth.AuthenticateAuthenticatorBeginRequest.create({
          authenticatorRef: getResourceRef(authn),
        })
      );

      if (response.challengeRequest?.type.oneofKind === `fido`) {
        try {
          console.log("Got req", response.challengeRequest.type.fido.request);
          const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(
            JSON.parse(response.challengeRequest.type.fido.request)
          );
          const credential = (await navigator.credentials.get({
            publicKey,
          })) as PublicKeyCredential;

          console.log("Got credential", credential.toJSON());

          return await c.authenticateWithAuthenticator(
            Auth.AuthenticateWithAuthenticatorRequest.create({
              authenticatorRef: getResourceRef(authn),
              challengeResponse: {
                type: {
                  oneofKind: "fido",
                  fido: {
                    response: JSON.stringify(credential.toJSON()),
                  },
                },
              },
            })
          );
        } catch (err) {
          console.log("fido get err", err);
          throw err;
        }
      }
    },
    onSuccess: (r) => {
      window.location.href = "/callback/success";
    },
    onError: (resp) => {},
  });

  React.useEffect(() => {
    mutation.mutate();
  }, []);

  return <div></div>;
};

const Authenticator = (props: { authn: Auth.Authenticator }) => {
  const { authn } = props;
  let [open, setOpen] = React.useState(false);
  let [isEdit, setIsEdit] = React.useState(false);
  let [displayName, setDisplayName] = React.useState<string>("");

  const isDelete = useDisclosure(false);
  const c = getClientAuth();

  const mutationDelete = useMutation({
    mutationFn: async () => {
      await c.deleteAuthenticator(
        DeleteOptions.create({
          uid: authn.metadata?.uid,
          name: authn.metadata?.name,
        })
      );
    },
    onSuccess: (r) => {
      queryClient.invalidateQueries({
        queryKey: ["getAvailableAuthenticator"],
      });
    },
    onError: (resp) => {},
  });

  const mutationUpdate = useMutation({
    mutationFn: async () => {
      let req = Auth.Authenticator.clone(authn);
      req.spec!.displayName = displayName ?? "";
      await c.updateAuthenticator(req);
    },
    onSuccess: (r) => {
      queryClient.invalidateQueries({
        queryKey: ["getAvailableAuthenticator"],
      });
    },
    onError: (resp) => {},
  });

  return (
    <div
      className={twMerge(
        "w-full",
        "p-4 mb-3",

        `bg-slate-50 hover:bg-white`,
        `border-[2px] border-gray-200 rounded-lg`,
        "w-full",

        "transition-all duration-300",
        // "hover:bg-transparent",
        "py-4 px-2",
        "font-semibold",
        "rounded-xl",
        "shadow-sm shadow-slate-200",
        "border-[2px] border-slate-300",
        "mb-4"
      )}
    >
      <div className="flex items-center">
        <div className="w-full font-bold flex-1">
          <div className="flex items-center">
            <span className="text-white bg-slate-800 rounded-lg p-2 text-xs shadow-md">
              {authn.status?.type === Auth.Authenticator_Status_Type.FIDO &&
                "FIDO"}
              {authn.status?.type === Auth.Authenticator_Status_Type.TOTP &&
                "TOTP"}
            </span>{" "}
            {!isEdit && (
              <span className="text-black ml-1">
                {authn.spec?.displayName ?? authn.metadata!.name}
              </span>
            )}
            {isEdit && (
              <div className="flex items-center">
                <Input
                  placeholder="My Authenticator"
                  variant={"unstyled"}
                  value={authn.spec?.displayName}
                  onChange={(e) => {
                    setDisplayName(e.target.value);
                  }}
                />

                <ActionIcon
                  aria-label="Submit"
                  onClick={() => {
                    mutationUpdate.mutate();
                  }}
                >
                  <IoMdSend />
                </ActionIcon>
              </div>
            )}
            <Tooltip label="Edit Display Name">
              <ActionIcon
                variant="transparent"
                aria-label="Edit Display Name"
                onClick={() => {
                  setIsEdit(!isEdit);
                }}
              >
                {isEdit ? (
                  <MdEditOff className="text-slate-500" />
                ) : (
                  <MdEdit className="text-slate-500" />
                )}
              </ActionIcon>
            </Tooltip>
          </div>
          <div className="text-xs mt-1 mb-2 text-slate-500">
            <span>Created </span>
            <TimeAgo rfc3339={authn.metadata?.createdAt} />
          </div>
        </div>
        <div className="flex flex-col items-center">
          <Button
            onClick={() => {
              setOpen(!open);
            }}
            variant={open ? "outline" : undefined}
          >
            {open ? "Cancel" : "Authenticate"}
          </Button>
          <Button
            className="mt-3 !rounded-md !transition-all !duration-500 !border-slate-500"
            fullWidth
            size="compact-xs"
            variant="outline"
            onClick={isDelete[1].open}
          >
            <span className="text-slate-600">Delete</span>
          </Button>
        </div>
      </div>

      <Collapse in={open} transitionDuration={300}>
        {open && (
          <div>
            {authn.status?.type === Auth.Authenticator_Status_Type.FIDO && (
              <Fido authn={authn} />
            )}

            {authn.status?.type === Auth.Authenticator_Status_Type.TOTP && (
              <TOTP authn={authn} />
            )}
          </div>
        )}
      </Collapse>

      <Modal opened={isDelete[0]} onClose={isDelete[1].close} centered>
        <div className="font-bold text-xl mb-4">
          {`Do you really want to delete this Authenticator?`}
        </div>

        <div className="mt-4 flex justify-end items-center">
          <Button variant="outline" onClick={isDelete[1].close}>
            Cancel
          </Button>
          <Button
            className="ml-4"
            loading={mutationDelete.isPending}
            onClick={() => {
              mutationDelete.mutate();
            }}
            autoFocus
          >
            Yes, Delete
          </Button>
        </div>
      </Modal>
    </div>
  );
};

const devList = Auth.AuthenticatorList.create({
  items: [
    Auth.Authenticator.create({
      metadata: {
        name: "totp-1",
        createdAt: Timestamp.now(),
      },
      spec: {
        displayName: "Google Authenticator",
      },
      status: {
        type: Auth.Authenticator_Status_Type.TOTP,
      },
    }),
    Auth.Authenticator.create({
      metadata: {
        name: "fido-1",
        createdAt: Timestamp.now(),
      },
      spec: {
        displayName: "My FIDO Key",
      },
      status: {
        type: Auth.Authenticator_Status_Type.FIDO,
      },
    }),
  ],
});

const AvailableAuthenticators = (props: {
  resp: Auth.GetAvailableAuthenticatorResponse;
}) => {
  const { resp } = props;

  return (
    <div className="w-full">
      {resp.availableAuthenticators.map((x) => (
        <Authenticator authn={x} />
      ))}
    </div>
  );
};

const Page = () => {
  const c = getClientAuth();

  const { isError, isLoading, data } = useQuery({
    queryKey: ["getAvailableAuthenticator"],
    queryFn: async () => {
      return await c.getAvailableAuthenticator({});
    },
  });

  if (isDev()) {
    return (
      <div>
        {devList.items.map((x) => (
          <Authenticator authn={x} />
        ))}
      </div>
    );
  }

  if (isLoading) {
    return <></>;
  }
  if (!data) {
    return <></>;
  }

  return (
    <div>
      <div className="container mx-auto mt-2 p-2 md:p-4 w-full max-w-lg">
        {data.response.mainAuthenticator && (
          <div>
            {data.response.mainAuthenticator.status?.type ===
              Auth.Authenticator_Status_Type.FIDO && (
              <Fido authn={data.response.mainAuthenticator} />
            )}

            {data.response.mainAuthenticator.status?.type ===
              Auth.Authenticator_Status_Type.TOTP && (
              <TOTP authn={data.response.mainAuthenticator} />
            )}
          </div>
        )}

        <AvailableAuthenticators resp={data.response} />
      </div>
    </div>
  );
};

export default Page;
