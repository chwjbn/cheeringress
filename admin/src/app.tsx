import type { Settings as LayoutSettings } from '@ant-design/pro-layout';
import { PageLoading } from '@ant-design/pro-layout';
import type { RequestConfig, RunTimeLayoutConfig } from 'umi';
import { history } from 'umi';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import { message } from 'antd';
import { clearLoginAccessToken, getLoginAccessToken } from './services/ant-design-pro/api';

const loginPath = '/account/login';

const noNeedAuthPath = ['/account/login', '/account/reg', '/account/find'];

/** 获取用户信息比较慢的时候会展示一个 loading */
export const initialStateConfig = {
  loading: <PageLoading />,
};

const isNoNeedAuthPath = (): boolean => {
  const { location } = history;

  let noNeedAuthFlag = false;

  noNeedAuthPath.forEach(function (dv) {
    if (dv === location.pathname) {
      noNeedAuthFlag = true;
    }
  });

  return noNeedAuthFlag;
};


/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.AccessTokenData;
  fetchUserInfo?: () => Promise<API.AccessTokenData | undefined>;
}> {

  const fetchUserInfo = async () => {
    try {
      const xData = getLoginAccessToken();
      return xData;
    } catch (error) {
      if (!isNoNeedAuthPath()) {
        history.push(loginPath);
      }
    }
    return undefined;
  };

  if (!isNoNeedAuthPath()) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: {},
    };
  }


  return {
    fetchUserInfo,
    settings: {},
  };

  
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState }) => {
  return {
    rightContentRender: () => <RightContent />,
    disableContentMargin: false,
    footerRender: () => <Footer />,
    onPageChange: () => {
      if (!isNoNeedAuthPath()) {
        // 如果没有登录，重定向到 login
        if (!initialState?.currentUser?.username) {
          history.push(loginPath);
        }
      }
    },
    links: [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    ...initialState?.settings,
  };
};

// 自动添加用户身份信息
const authHeaderInterceptor = (url: string, options: any) => {
  const authHeader = { Authorization: '' };

  const xTenantUserInfo = getLoginAccessToken();

  if (xTenantUserInfo) {
    authHeader.Authorization = `${xTenantUserInfo.data_id}`;
  }

  return {
    url: `${url}`,
    options: { ...options, interceptors: true, headers: authHeader },
  };
};

const apiErrorResponseInterceptor = async (response: Response) => {
  
  const xResp = response.clone();

  if (!xResp.ok) {
    message.error('网络异常['+xResp.statusText+']');
  }

  const xJsonData=await xResp.json();

  if(xJsonData&&xJsonData.error_code){

    if(xJsonData.error_code=='401'){
      clearLoginAccessToken();
      window.location.href=loginPath;
    }

  }

  return response;
};

const apiErrorProcess = (error: any) => {
  window.console.error('apiErrorProcess', error);
};

export const request: RequestConfig = {
  timeout: 30 * 1000,
  errorHandler: apiErrorProcess,
  requestInterceptors: [authHeaderInterceptor],
  responseInterceptors: [apiErrorResponseInterceptor],
};
