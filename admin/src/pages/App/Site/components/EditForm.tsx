import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import { ProFormDigit } from '@ant-design/pro-form';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressSiteInfo,
  CtlIngressSiteSave,
  CtlIngressNamespaceMap,
  CtlIngressActionBackendMap,
  CtlIngressActionStaticMap,
  DataMapAppDataIngressActionType,
  DataMapAppDataIngressSiteAuthNeed,
  DataMapRuleStringOp,
} from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Spin } from 'antd';
import { Access } from 'umi';

export type EditFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  id?: string;
};

const EditForm: React.FC<EditFormProps> = (props) => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);

  const xFormRef = useRef<ProFormInstance<API.AppDataIngressSite>>();
  const [xNeedAuth, setNeedAuth] = useState<boolean>(false);
  const [xIngressActionValueMapData, setIngressActionValueMapData] = useState<any>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadNamespaceMapData = async () => {
    const dataList: any[] = [];

    const xReq = {};

    const xRespData = await CtlIngressNamespaceMap(xReq);

    if (xRespData && xRespData.data) {
      xRespData.data.forEach((xItem) => {
        dataList.push({ label: xItem.data_name, value: xItem.data_id });
      });
    }

    return dataList;
  };

  const loadIngressActionValueMapData = async () => {
    const xDataList: any[] = [];

    const xActionType = xFormRef.current?.getFieldValue('action_type');

    const xNamespaceId = xFormRef.current?.getFieldValue('namespace_id');
    const xReq = { data_id: xNamespaceId };

    if (!xNamespaceId) {
      return;
    }

    if (xActionType == 'backend') {
      const xRespData = await CtlIngressActionBackendMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    if (xActionType == 'static') {
      const xRespData = await CtlIngressActionStaticMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    setIngressActionValueMapData(xDataList);
  };

  const loadAuthFields = () => {
    const xNeedAuthFlag = xFormRef.current?.getFieldValue('auth_need') == 'yes';
    setNeedAuth(xNeedAuthFlag);
  };

  const loadDataInfo = async (id?: string) => {
    setIsLoading(true);

    if (id) {
      const respData = await CtlIngressSiteInfo({ data_id: id });
      if (respData && respData.error_code === '0') {
        if (respData.data) {
          const xDataInfo = respData.data;
          xFormRef.current?.setFieldsValue(xDataInfo);
          loadAuthFields();
          await loadIngressActionValueMapData();
        }
      }
    }

    setIsLoading(false);
  };

  // 加载
  useEffect(() => {
    if (props.modalVisible) {
      loadDataInfo(props.id);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.id, props.modalVisible]);

  const onSaveData = async (xData: API.AppDataIngressSite) => {
    xData.data_id = props.id;

    const xResp = await CtlIngressSiteSave(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressSite>
      title={'编辑站点'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressSite);
      }}
      formRef={xFormRef}
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProFormSelect
          name="namespace_id"
          label="网关空间"
          placeholder="请选择网关空间"
          rules={[{ required: true, message: '请选择网关空间.' }]}
          request={loadNamespaceMapData}
          fieldProps={{
            showSearch: true,
            onChange: (value) => {
              xFormRef.current?.resetFields();
              xFormRef.current?.setFieldsValue({ namespace_id: value });
            },
          }}
        />

        <ProFormText
          name="title"
          label="站点名称"
          placeholder="请输入站点名称"
          rules={[{ required: true, message: '请输入站点名称.' }]}
        />

        <ProFormDigit
          name="order_no"
          label="站点序号"
          placeholder="请输入站点序号"
          rules={[{ required: true, message: '请输入站点序号.' }]}
        />

        <ProFormSelect
          name="match_op"
          label="域名匹配方式"
          valueEnum={DataMapRuleStringOp}
          placeholder="请选择域名匹配方式"
          rules={[{ required: true, message: '请选择域名匹配方式.' }]}
        />

        <ProFormText
          name="match_value"
          label="域名匹配内容"
          placeholder="请输入域名匹配内容"
          rules={[{ required: true, message: '请输入域名匹配内容.' }]}
        />

        <ProFormSelect
          name="auth_need"
          label="启用访问认证"
          valueEnum={DataMapAppDataIngressSiteAuthNeed}
          placeholder="请选择是否启用访问认证"
          rules={[{ required: true, message: '请选择是否启用访问认证.' }]}
          fieldProps={{
            onChange() {
              loadAuthFields();
            },
          }}
        />

        <Access key={'auth_need_data'} accessible={xNeedAuth}>
          <ProFormText
            name="auth_user_name"
            label="访问认证用户名"
            placeholder="请输入访问认证用户名"
            rules={[{ required: false, message: '请输入访问认证用户名.' }]}
          />

          <ProFormText
            name="auth_password"
            label="访问认证密码"
            placeholder="请输入访问认证密码"
            rules={[{ required: false, message: '请输入访问认证密码.' }]}
          />
        </Access>

        <ProFormSelect
          name="action_type"
          label="响应类型"
          valueEnum={DataMapAppDataIngressActionType}
          placeholder="请选择响应类型"
          rules={[{ required: true, message: '请选择响应类型.' }]}
          fieldProps={{
            onChange() {
              const xFormDataVal = xFormRef.current?.getFieldsValue();
              if (xFormDataVal) {
                xFormDataVal.action_value = undefined;
                xFormRef.current?.setFieldsValue(xFormDataVal);
              }
              loadIngressActionValueMapData();
            },
          }}
        />

        <ProFormSelect
          name="action_value"
          label="响应内容"
          placeholder="请选择响应内容"
          fieldProps={{ options: xIngressActionValueMapData, showSearch: true }}
          rules={[{ required: true, message: '请选择响应内容.' }]}
        />
        <ProFormSelect
          name="state"
          label="状态"
          valueEnum={{
            enable: { text: '启用', status: 'enable' },
            disable: { text: '禁用', status: 'disable' },
          }}
          placeholder="请选择状态"
          rules={[{ required: true, message: '请选择状态.' }]}
        />
      </Spin>
    </ModalForm>
  );
};

export default EditForm;
