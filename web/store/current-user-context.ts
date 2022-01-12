import * as React from "react";

export interface CurrentUser {
  uuid: string | null;
  email: string | null;
  displayName: string | null;
  role: string | null;
  balance: number | null;
}

export interface CurrentUserAppContextInterface {
  uuid: string | null;
  email: string | null;
  displayName: string | null;
  role: string | null;
  balance: number | null;
  fetchCurrentUser: (currentUser: CurrentUser) => void;
  removeCurrentUser: () => void;
}

const CurrentUserAppCtx = React.createContext<CurrentUserAppContextInterface | null>(null);

export default CurrentUserAppCtx;