import * as React from "react";

export interface CurrentUser {
  uuid: string | null;
  email: string | null;
  displayName: string | null;
  role: string | null;
}

export interface CurrentUserAppContextInterface {
  uuid: string | null;
  email: string | null;
  displayName: string | null;
  role: string | null;
  fetchCurrentUser: (currentUser: CurrentUser) => void;
  removeCurrentUser: () => void;
}

const CurrentUserAppCtx = React.createContext<CurrentUserAppContextInterface | null>(null);

export default CurrentUserAppCtx;