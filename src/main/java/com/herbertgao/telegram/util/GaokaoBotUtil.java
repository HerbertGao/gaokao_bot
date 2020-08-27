package com.herbertgao.telegram.util;

import com.herbertgao.telegram.bot.TemplateReplace;
import com.herbertgao.telegram.database.entity.ExamDate;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

/**
 * @program: gaokao_bot
 * @description: CommonService
 * @author: HerbertGao
 * @create: 2019-06-09 01:52
 **/
public class GaokaoBotUtil {

    /**
     * 是考试时间
     *
     * @param now
     * @return
     */
    public static Boolean isExamTime(ExamDate exam, LocalDateTime now) {
        LocalDateTime examBeginDate = exam.getExamBeginDate();
        return examBeginDate.isBefore(now) || examBeginDate.isEqual(now);
    }

    /**
     * 获取倒计时文字
     *
     * @param exam
     * @param now
     * @param template
     * @return
     */
    public static String getCountDownString(ExamDate exam, LocalDateTime now, String template) {
        if (isExamTime(exam, now)) {
            return exam.getExamDesc() + "正在进行中！" + System.getProperty("line.separator");
        } else {
            String rtn = template;
            String[] searchList = {
                    TemplateReplace.EXAM_YEAR,
                    TemplateReplace.EXAM,
                    TemplateReplace.EXAM_S,
                    TemplateReplace.TIME
            };
            String[] replacementList = {
                    exam.getExamYear().toString(),
                    exam.getExamDesc(),
                    exam.getShortDesc(),
                    getCountDownTime(exam, now)
            };
            rtn = StringUtils.replaceEach(rtn, searchList, replacementList);
            rtn += System.getProperty("line.separator");
            return rtn;
        }
    }

    /**
     * 获取倒计时时间文字
     *
     * @param exam
     * @param now
     * @return
     */
    public static String getCountDownTime(ExamDate exam, LocalDateTime now) {

        StringBuilder rtn = new StringBuilder();

        LocalDateTime examBeginDate = exam.getExamBeginDate();
        long daysDiff = ChronoUnit.DAYS.between(now, examBeginDate);
        long hoursDiff = ChronoUnit.HOURS.between(now, examBeginDate) - daysDiff * 24;
        long minutesDiff = ChronoUnit.MINUTES.between(now, examBeginDate) - daysDiff * 24 * 60 - hoursDiff * 60;
        long secondsDiff = ChronoUnit.SECONDS.between(now, examBeginDate) - daysDiff * 24 * 60 * 60 - hoursDiff * 60 * 60 - minutesDiff * 60;

        if (daysDiff > 0) {
            rtn.append(daysDiff).append("天");
        }
        if (hoursDiff > 0) {
            rtn.append(hoursDiff).append("小时");
        }
        if (minutesDiff > 0) {
            rtn.append(minutesDiff).append("分钟");
        }
        if (secondsDiff > 0) {
            rtn.append(secondsDiff).append("秒");
        }

        return rtn.toString();
    }

}
