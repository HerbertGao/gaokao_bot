package com.herbertgao.telegram.database.service;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.mapper.UserTemplateMapper;
import com.herbertgao.telegram.util.IdUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @program: gaokao_bot
 * @description:
 * @author: HerbertGao
 * @create: 2020-04-08 18:13
 **/
@Service
public class UserTemplateService {

    @Autowired
    private UserTemplateMapper mapper;

    /**
     * 通过用户ID获取模板列表
     *
     * @param userId
     * @return
     */
    public List<UserTemplate> getUserTemplateListByUserId(Long userId) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getUserId, userId);
        return mapper.selectList(wrapper);
    }

    /**
     * @param userId
     * @param examUsername
     */
    public Long insertDefaultTemplateWithUsername(Long userId, String examUsername) {
        UserTemplate defaultTemplate = getDefaultTemplate();
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName("@" + examUsername);
        userTemplate.setTemplateContent(defaultTemplate.getTemplateContent() + "\n@" + examUsername);
        mapper.insert(userTemplate);
        return id;
    }

    /**
     * @param userId
     * @param template
     */
    public Long insertTemplateWithUsername(Long userId, String template, String templateName) {
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName(templateName);
        userTemplate.setTemplateContent(template);
        mapper.insert(userTemplate);
        return id;
    }

    /**
     * 获取默认模板列表
     *
     * @return
     */
    public UserTemplate getDefaultTemplate() {
        return getUserTemplateListByUserId(0L).get(0);
    }

    /**
     * @param userId
     * @param id
     * @return
     */
    public UserTemplate getTemplateById(Long userId, String id) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getId, Long.parseLong(id))
                .eq(UserTemplate::getUserId, userId);
        return mapper.selectOne(wrapper);
    }

    /**
     * @param userId
     * @param id
     */
    public void remove(Long userId, String id) {
        QueryWrapper<UserTemplate> wrapper = new QueryWrapper<>();
        wrapper.lambda()
                .eq(UserTemplate::getId, Long.parseLong(id))
                .eq(UserTemplate::getUserId, userId);
        mapper.delete(wrapper);
    }

    /**
     * @param template
     */
    public void update(UserTemplate template) {
        mapper.updateById(template);
    }
}
