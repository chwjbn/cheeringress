/**
 * @see https://umijs.org/zh-CN/plugins/plugin-access
 * */

export type AppAccess = {
  canAdmin: boolean;
  canDev: boolean;
  canOp: boolean;
  canQc: boolean;
  canProd: boolean;
  isCurrentUser: (userId?: number) => boolean;
  isCurrentRole: (userRole?: string) => boolean;
};

export default function access(initialState: {
  currentUser?: API.AccessTokenData | undefined;
}): AppAccess {
  const { currentUser } = initialState || {};

  return {
    canAdmin: currentUser && currentUser.role === 'admin' ? true : false,
    canDev: currentUser && currentUser.role === 'dev' ? true : false,
    canOp: currentUser && currentUser.role === 'op' ? true : false,
    canQc: currentUser && currentUser.role === 'qc' ? true : false,
    canProd: currentUser && currentUser.role === 'prod' ? true : false,
    isCurrentUser: (userId?: number)=>{ return currentUser?.account_id==userId},
    isCurrentRole: (userRole?: string)=>{ return currentUser?.role==userRole},
  };
}
