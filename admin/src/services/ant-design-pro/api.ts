import { request } from 'umi';

import { ClearLocalCache, GetLocalCache, SetLocalCache } from './lib';

export function setLoginAccessToken(loginAccessToken?: API.AccessTokenData) {
  SetLocalCache('x_access_token', JSON.stringify(loginAccessToken));
}

export function getLoginAccessToken(): undefined | API.AccessTokenData {
  let xLoginAccessToken: undefined | API.AccessTokenData = undefined;

  const xLoginAccessTokenJson = GetLocalCache('x_access_token');
  if (!xLoginAccessTokenJson) {
    return xLoginAccessToken;
  }

  xLoginAccessToken = JSON.parse(xLoginAccessTokenJson);

  return xLoginAccessToken;
}

export function clearLoginAccessToken() {
  ClearLocalCache();
}

export function getLoginTokenId() {
  const xLoginAccessToken = getLoginAccessToken();

  if (xLoginAccessToken) {
    return xLoginAccessToken.data_id;
  }

  return '';
}

export async function CtlCheckCodeImage() {
  return request<API.AccountCheckCodeImageResponse>('/xapi/account/check-code-image', {
    method: 'GET',
  });
}

export async function GetCheckCodeImage(): Promise<string> {
  let xData = '';

  const xResp = await CtlCheckCodeImage();

  if (xResp && xResp.error_code && xResp.error_code == '0') {
    if (xResp.data?.code_id && xResp.data?.image_data) {
      xData = xResp.data?.image_data;
      SetLocalCache('x_check_code_image_id', xResp.data?.code_id);
    }
  }

  return xData;
}

export function GetCheckCodeImageId(): string {
  return GetLocalCache('x_check_code_image_id');
}

export async function CtlAccountLogin(reqData: API.AccountLoginRequest) {
  return request<API.AccountLoginResponse>('/xapi/account/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlUserGetCurrent() {
  return request<API.UserGetInfoResponse>('/xapi/user/get-current', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: {},
  });
}

export async function CtlUserUpdateInfo(reqData: API.UserInfoUpdateRequest) {
  return request<API.BaseDataResponse>('/xapi/user/update-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlUserUpdatePassword(reqData: API.UserSecurityUpdateRequest) {
  return request<API.BaseDataResponse>('/xapi/user/update-security', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespacePageData(reqData: API.AppDataIngressNamespacePageRequest) {
  const xDataRet = request<API.AppDataIngressNamespacePageResponse>(
    '/xapi/ingress/namespace-page',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: reqData,
    },
  );

  return xDataRet;
}

export async function CtlIngressNamespaceInfo(reqData: API.AppDataIdRequest) {
  return request<API.AppDataIngressNamespaceInfoResponse>('/xapi/ingress/namespace-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespaceAdd(reqData: API.AppDataIngressNamespaceAddRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/namespace-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespaceSave(reqData: API.AppDataIngressNamespaceSaveRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/namespace-save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespaceRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/namespace-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespacePublish(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/namespace-publish', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressNamespaceMap(reqData: API.AppDataIngressNamespaceMapRequest) {
  return request<API.AppDataIngressNamespaceMapResponse>('/xapi/ingress/namespace-map', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendPageData(
  reqData: API.AppDataIngressActionBackendPageRequest,
) {
  const xDataRet = request<API.AppDataIngressActionBackendPageResponse>(
    '/xapi/ingress/action-backend-page',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: reqData,
    },
  );

  return xDataRet;
}

export async function CtlIngressActionBackendInfo(reqData: API.AppDataIdRequest) {
  return request<API.AppDataIngressActionBackendInfoResponse>('/xapi/ingress/action-backend-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendAdd(
  reqData: API.AppDataIngressActionBackendAddRequest,
) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-backend-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendSave(
  reqData: API.AppDataIngressActionBackendSaveRequest,
) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-backend-save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-backend-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendMap(
  reqData: API.AppDataIngressActionBackendMapRequest,
) {
  return request<API.AppDataIngressActionBackendMapResponse>('/xapi/ingress/action-backend-map', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendNodePage(
  reqData: API.AppDataIngressActionBackendNodePageRequest,
) {
  return request<API.AppDataIngressActionBackendNodePageResponse>(
    '/xapi/ingress/action-backend-node-page',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: reqData,
    },
  );
}

export async function CtlIngressActionBackendNodeAdd(
  reqData: API.AppDataIngressActionBackendNodeAddRequest,
) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-backend-node-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionBackendNodeRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-backend-node-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionStaticPageData(
  reqData: API.AppDataIngressActionStaticPageRequest,
) {
  const xDataRet = request<API.AppDataIngressActionStaticPageResponse>(
    '/xapi/ingress/action-static-page',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: reqData,
    },
  );

  return xDataRet;
}

export async function CtlIngressActionStaticInfo(reqData: API.AppDataIdRequest) {
  return request<API.AppDataIngressActionStaticInfoResponse>('/xapi/ingress/action-static-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionStaticAdd(reqData: API.AppDataIngressActionStaticAddRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-static-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionStaticSave(
  reqData: API.AppDataIngressActionStaticSaveRequest,
) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-static-save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressActionStaticRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/action-static-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}


export async function CtlIngressActionStaticMap(reqData: API.AppDataIngressActionStaticMapRequest) {
  return request<API.AppDataIngressActionStaticMapResponse>('/xapi/ingress/action-static-map', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}


export async function CtlIngressSitePageData(reqData: API.AppDataIngressSitePageRequest) {
  const xDataRet = request<API.AppDataIngressSitePageResponse>('/xapi/ingress/site-page', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });

  return xDataRet;
}

export async function CtlIngressSiteInfo(reqData: API.AppDataIdRequest) {
  return request<API.AppDataIngressSiteInfoResponse>('/xapi/ingress/site-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteAdd(reqData: API.AppDataIngressSiteAddRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/site-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteSave(reqData: API.AppDataIngressSiteSaveRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/site-save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/site-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteRulePage(reqData: API.AppDataIngressSiteRulePageRequest) {
  return request<API.AppDataIngressSiteRulePageResponse>('/xapi/ingress/site-rule-page', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteRuleAdd(reqData: API.AppDataIngressSiteRuleAddRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/site-rule-add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}


export async function CtlIngressSiteRuleInfo(reqData: API.AppDataIdRequest) {
  return request<API.AppDataIngressSiteInfoResponse>('/xapi/ingress/site-rule-info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlIngressSiteRuleRemove(reqData: API.AppDataIdRequest) {
  return request<API.BaseDataResponse>('/xapi/ingress/site-rule-remove', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

export async function CtlAliyunOSSUploadArg(reqData: any) {
  return request<API.AliyunOSSUploadArgResponse>('/xapi/api/upload-arg', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: reqData,
  });
}

const DataMapAppDataState = {
  enable: { text: '??????', status: 'enable' },
  disable: { text: '??????', status: 'disable' },
};

const DataMapAppDataIngressActionBackendBalanceType = {
  IPHash: { text: 'IPHash??????', status: 'IPHash' },
  RoundRobin: { text: '????????????', status: 'RoundRobin' },
};

const DataMapAppDataIngressActionStaticContentType = [
  { label: 'text/plain', value: 'text/plain; charset=utf-8' },
  { label: 'text/html', value: 'text/html; charset=utf-8' },
  { label: 'text/xml', value: 'text/xml; charset=utf-8' },
  { label: 'application/json', value: 'application/json; charset=utf-8' },
  { label: 'application/octet-stream', value: 'application/octet-stream; charset=utf-8' },
  { label: 'image/png', value: 'image/png; charset=utf-8' },
  { label: 'image/jpeg', value: 'image/jpeg; charset=utf-8' },
  { label: 'image/x-icon', value: 'image/x-icon; charset=utf-8' },
];

const DataMapAppDataIngressActionStaticDataType = {
  PlainText: { text: '??????????????????', status: 'PlainText' },
  Base64Data: { text: 'BASE64????????????', status: 'Base64Data' },
  HttpResContent: { text: 'HTTP????????????', status: 'HttpResContent' },
  HttpResZip: { text: 'HTTP???????????????', status: 'HttpResZip' },
  Http301Redirect: { text: 'HTTP?????????301', status: 'Http301Redirect' },
};

const DataMapAppDataIngressSiteAuthNeed = {
  yes: { text: '???', status: 'yes' },
  no: { text: '???', status: 'no' },
};


const DataMapAppDataIngressActionType = {
  backend: { text: '????????????', status: 'backend' },
  static: { text: '????????????', status: 'static' },
};

const DataMapHttpMethod = {
  ALL: { text: '??????', status: 'ALL' },
  GET: { text: 'GET', status: 'GET' },
  POST: { text: 'POST', status: 'POST' },
  HEAD: { text: 'HEAD', status: 'HEAD' },
  OPTIONS: { text: 'OPTIONS', status: 'OPTIONS' },
  PUT: { text: 'PUT', status: 'PUT' },
  DELETE: { text: 'DELETE', status: 'DELETE' },
};

const DataMapHttpTargetItem = {
  s_header_url: { text: '???????????????URL', status: 's_header_url' },
  s_header_referer: { text: '???????????????URL', status: 's_header_referer' },
  s_header_useragent: { text: '????????????????????????', status: 's_header_useragent' },
  s_header_cookie: { text: '?????????COOKIE', status: 's_header_cookie' },
  s_header_x_tenant: { text: '???????????????[X-Tenant]???', status: 's_header_x_tenant' },
  s_client_ip: { text: '?????????IP??????', status: 's_client_ip' },
};

const DataMapRuleStringOp = {
  eq: { text: '??????', status: 'eq' },
  contain: { text: '????????????', status: 'contain' },
  regex: { text: '????????????', status: 'regex' },
  in: { text: '???????????????', status: 'in' },
  neq: { text: '?????????', status: 'neq' },
  notcontain: { text: '???????????????', status: 'notcontain' },
  notregex: { text: '???????????????', status: 'notregex' },
  notin: { text: '??????????????????', status: 'notin' },
};

export {
  DataMapAppDataState,
  DataMapAppDataIngressActionBackendBalanceType,
  DataMapAppDataIngressActionStaticContentType,
  DataMapAppDataIngressActionStaticDataType,
  DataMapAppDataIngressSiteAuthNeed,
  DataMapAppDataIngressActionType,
  DataMapHttpMethod,
  DataMapHttpTargetItem,
  DataMapRuleStringOp,
};
